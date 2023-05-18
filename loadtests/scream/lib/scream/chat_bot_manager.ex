defmodule Scream.ChatBotManager do
  use DynamicSupervisor

  def init(:ok) do
    DynamicSupervisor.init(strategy: :one_for_one)
  end

  def new_bot(opts) do
    DynamicSupervisor.start_child(__MODULE__, {Scream.ChatBot, opts})
  end
end
