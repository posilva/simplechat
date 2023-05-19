defmodule Scream.ChatBot do
  use GenServer

  require Logger
  require Mint.HTTP

  defstruct [
    :conn,
    :websocket,
    :request_ref,
    :caller,
    :status,
    :resp_headers,
    :bot_id,
    :bot_room,
    :closing?
  ]

  def load(bot) do
    GenServer.call(bot, :start_load)
  end

  def send_message(pid, text) do
    GenServer.call(pid, {:send_text, text})
  end

  def start_link(initial) do
    GenServer.start_link(__MODULE__, initial)
  end

  @impl GenServer
  @spec init(any) ::
          {:ok,
           %Scream.ChatBot{
             caller: nil,
             closing?: nil,
             conn: nil,
             request_ref: nil,
             resp_headers: nil,
             status: nil,
             websocket: nil
           }}
  def init(url) do
    idx = Enum.random([1, 2, 3, 4])
    room = "room_#{idx}"
    from = UUID.uuid4()

    complete_url = "#{url}?room=#{room}&id=#{from}"
    GenServer.cast(self(), {:connect, complete_url})

    state = %__MODULE__{}
    state = put_in(state.bot_id, from)
    state = put_in(state.bot_room, room)
    Logger.info("init bot: #{inspect(state.bot_id)} connected to room #{inspect(state.bot_room)}")
    {:ok, state}
  end

  @impl GenServer
  def handle_call(:start_load, _from, state) do
    send_load_msg(state)
    {:reply, :ok, state}
  end

  @impl GenServer
  def handle_call({:send_text, text}, _from, state) do
    {:ok, state} = send_frame(state, {:text, text})
    {:reply, :ok, state}
  end

  @impl GenServer
  def handle_cast({:connect, url}, state) do
    uri = URI.parse(url)

    http_scheme =
      case uri.scheme do
        "ws" -> :http
        "wss" -> :https
      end

    ws_scheme =
      case uri.scheme do
        "ws" -> :ws
        "wss" -> :wss
      end

    path =
      case uri.query do
        nil -> uri.path
        query -> uri.path <> "?" <> query
      end

    with {:ok, conn} <- Mint.HTTP.connect(http_scheme, uri.host, uri.port),
         {:ok, conn, ref} <- Mint.WebSocket.upgrade(ws_scheme, conn, path, []) do
      state = %{state | conn: conn, request_ref: ref}
      {:noreply, state}
    else
      {:error, reason} ->
        Logger.error("failed to connect: #{inspect(reason)}")
        {:noreply, state}

      {:error, conn, reason} ->
        Logger.error("failed to connect: #{inspect(reason)}")
        {:noreply, put_in(state.conn, conn)}
    end
  end

  @impl GenServer
  def handle_info({:chat_msg, chat_msg}, state) do
    {:ok, state} = send_frame(state, {:text, chat_msg})
    send_load_msg(state)
    {:noreply, state}
  end

  @impl GenServer
  def handle_info(message, state) do
    case Mint.WebSocket.stream(state.conn, message) do
      {:ok, conn, responses} ->
        state = put_in(state.conn, conn) |> handle_responses(responses)
        if state.closing?, do: do_close(state), else: {:noreply, state}

      {:error, conn, reason, _responses} ->
        state = put_in(state.conn, conn) |> reply({:error, reason})
        {:noreply, state}

      :unknown ->
        {:noreply, state}
    end
  end

  defp handle_responses(state, responses)

  defp handle_responses(%{request_ref: ref} = state, [{:status, ref, status} | rest]) do
    put_in(state.status, status)
    |> handle_responses(rest)
  end

  defp handle_responses(%{request_ref: ref} = state, [{:headers, ref, resp_headers} | rest]) do
    put_in(state.resp_headers, resp_headers)
    |> handle_responses(rest)
  end

  defp handle_responses(%{request_ref: ref} = state, [{:done, ref} | rest]) do
    case Mint.WebSocket.new(state.conn, ref, state.status, state.resp_headers) do
      {:ok, conn, websocket} ->
        %{state | conn: conn, websocket: websocket, status: nil, resp_headers: nil}
        |> reply({:ok, :connected})
        |> handle_responses(rest)

      {:error, conn, reason} ->
        put_in(state.conn, conn)
        |> reply({:error, reason})
    end
  end

  defp handle_responses(%{request_ref: ref, websocket: websocket} = state, [
         {:data, ref, data} | rest
       ])
       when websocket != nil do
    case Mint.WebSocket.decode(websocket, data) do
      {:ok, websocket, frames} ->
        put_in(state.websocket, websocket)
        |> handle_frames(frames)
        |> handle_responses(rest)

      {:error, websocket, reason} ->
        Logger.error("failed to decode message #{inspect(reason)}")

        put_in(state.websocket, websocket)
        |> reply({:error, reason})
    end
  end

  defp handle_responses(state, [_response | rest]) do
    handle_responses(state, rest)
  end

  defp handle_responses(state, []), do: state

  defp send_frame(state, frame) do
    with {:ok, websocket, data} <- Mint.WebSocket.encode(state.websocket, frame),
         state = put_in(state.websocket, websocket),
         {:ok, conn} <- Mint.WebSocket.stream_request_body(state.conn, state.request_ref, data) do
      {:ok, put_in(state.conn, conn)}
    else
      {:error, %Mint.WebSocket{} = websocket, reason} ->
        {:error, put_in(state.websocket, websocket), reason}

      {:error, conn, reason} ->
        {:error, put_in(state.conn, conn), reason}
    end
  end

  def handle_frames(state, frames) do
    Enum.reduce(frames, state, fn
      # reply to pings with pongs
      {:ping, data}, state ->
        {:ok, state} = send_frame(state, {:pong, data})
        state

      {:close, _code, reason}, state ->
        Logger.debug("Closing connection: #{inspect(reason)}")
        %{state | closing?: true}

      {:text, _text}, state ->
        #Logger.debug("Received: #{inspect(text)}, sending back the reverse")
        # {:ok, state} = send_frame(state, {:text, text})
        state

      frame, state ->
        Logger.debug("Unexpected frame received: #{inspect(frame)}")
        state
    end)
  end

  defp do_close(state) do
    # Streaming a close frame may fail if the server has already closed
    # for writing.
    _ = send_frame(state, :close)
    Mint.HTTP.close(state.conn)
    {:stop, :normal, state}
  end

  defp reply(state, response) do
    if state.caller, do: GenServer.reply(state.caller, response)
    put_in(state.caller, nil)
  end

  defp send_load_msg(state) do
    chat_msg =
      Jason.encode!(%{
        From: state.bot_id,
        Payload: "message from #{state.bot_id} to room: #{state.bot_room}",
        To: state.bot_room
      })

    t = Enum.random(100..3000)
    Process.send_after(self(), {:chat_msg, chat_msg}, t)
  end
end
