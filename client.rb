#!/usr/bin/env ruby
require "xmlrpc/client"
require "json"

server = XMLRPC::Client.new("127.0.0.1", "/RPC2", 3000)
begin
  puts JSON.pretty_generate(server.call("Test.echo", "hogehoge"))

  puts JSON.pretty_generate(server.call("Test.sample", "hogehoge", 1, true, 3.14, ["hoge","fuga"], {
    "name"=>"hoge",
  }))
rescue XMLRPC::FaultException => e
  puts "Error:"
  puts e.faultCode
  puts e.faultString
end
