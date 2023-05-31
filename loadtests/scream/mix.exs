defmodule Scream.MixProject do
  use Mix.Project

  def project do
    [
      app: :scream,
      version: "0.1.0",
      elixir: "~> 1.14",
      start_permanent: Mix.env() == :prod,
      deps: deps(),
      releases: releases()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger, :crypto],
      mod: {Scream.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:elixir_uuid, "~> 1.2"},
      {:jason, "~> 1.4"},
      {:mint_web_socket, "~> 1.0"}
      # {:dep_from_hexpm, "~> 0.3.0"},
      # {:dep_from_git, git: "https://github.com/elixir-lang/my_dep.git", tag: "0.1.0"}
    ]
  end

  defp releases() do
    [
      scream: [
        include_executables_for: [:unix],
        applications: [runtime_tools: :permanent, crypto: :permanent],
        include_erts: false
      ]
    ]
  end
end
