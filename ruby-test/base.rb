# -*- encoding : utf-8 -*-
require 'json'
require 'ffi'

class Base
  extend FFI::Library
  ffi_lib './word-parser.so'

  # Cloud Methods
  attach_function :golangParseWords, [:pointer, :pointer], :string

  # Tree Methods
  attach_function :golangPostfixSets, [:pointer, :pointer], :string

end
