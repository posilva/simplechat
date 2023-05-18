defmodule ScreamTest do
  use ExUnit.Case
  doctest Scream

  test "greets the world" do
    assert Scream.hello() == :world
  end
end
