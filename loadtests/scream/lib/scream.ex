defmodule Scream do
  @moduledoc """
  Documentation for `Scream`.
  """

  @doc """
  Hello world.

  ## Examples

      iex> Scream.hello()
      :world

  """
def connect() do
      Scream.ChatBotManager.new_bot("ws://localhost:8081/ws?room=room1&id=1")
end

def load(n) do
  for _ <- 1..n do
      {:ok, bot} = Scream.ChatBotManager.new_bot("ws://localhost:8081/ws")
      Scream.ChatBot.load(bot)
  end
end
end
