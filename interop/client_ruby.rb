#!/usr/bin/env ruby
require "json"
require "xmlrpc/client"

server = XMLRPC::Client.new("127.0.0.1", "/RPC2", 3000)

begin
  echo = server.call("Test.echo", "hello from ruby")
  raise "echo mismatch: #{echo.inspect}" unless echo == "hello from ruby"

  payload = server.call(
    "Test.sample",
    "edge-case",
    -2_147_483_648,
    true,
    3.141592653589793,
    ["alpha", "beta"],
    { "name" => "ruby" }
  )

  reply = payload["Reply"] || payload
  raise "Arg1 mismatch: #{reply.inspect}" unless reply["Arg1"] == "edge-case"
  raise "Arg2 mismatch: #{reply.inspect}" unless reply["Arg2"] == -2_147_483_648
  raise "Arg3 mismatch: #{reply.inspect}" unless reply["Arg3"] == true
  raise "Arg4 mismatch: #{reply.inspect}" unless (reply["Arg4"] - 3.141592653589793).abs < 1e-12
  raise "Arg5 mismatch: #{reply.inspect}" unless reply["Arg5"] == ["alpha", "beta"]
  raise "Arg6 mismatch: #{reply.inspect}" unless reply["Arg6"]["name"] == "sample"

  puts JSON.pretty_generate({
    "ruby_echo" => echo,
    "ruby_sample" => reply
  })
rescue XMLRPC::FaultException => e
  warn "Ruby XML-RPC fault: #{e.faultCode} #{e.faultString}"
  exit 1
rescue StandardError => e
  warn "Ruby XML-RPC check failed: #{e.message}"
  exit 1
end
