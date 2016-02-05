# -*- encoding : utf-8 -*-
require_relative './base'

class Tree < Base
  attr_accessor :target, :responses, :comments
  ffi_lib './word-parser.so'

  attach_function :postfixSets, [:pointer, :pointer], :string
  attach_function :responseIds, [:pointer, :pointer], :string

  def initialize(responses, comments, target)
    self.responses  = responses
    self.comments   = comments
    self.target     = target
  end

  def response_ids
    @response_ids ||= responseIds(json_sentences, target)
  end

  def postfix_sets
    @postfix_sets ||= postfixSets(json_sentences, target)
  end

  private

  # TEST
  def json_sentences
    "{\"ResponseID\":\"2021707\",\"Sentences\":[[\"I love when people make trees\"]]}{\"ResponseID\":\"2021708\",\"Sentences\":[[\"I also love eating pizza \"]]}{\"ResponseID\":\"2021709\",\"Sentences\":[[\"eating light and green is very good\"],[\"love me, tender... ;D\"]]}{\"ResponseID\":\"2021710\",\"Sentences\":[[\"I love you\"],[\"I love to eat pizza for breakfast\"],[\"I would like to know what love is all about\"]]}{\"ResponseID\":\"2021711\",\"Sentences\":[[\"my pizza is what I give you\"],[\"I eat pizza and nothing else\"],[\"but I never drink coca-cola\"]]}{\"ResponseID\":\"2021712\",\"Sentences\":[[\"I would love you to help me more !\"]]}{\"ResponseID\":\"2021713\",\"Sentences\":[[\"I love to eat with friends and family\"]]}"
  end

  # #  PRODUCTION
  # def json_sentences
  #   (responses + comments).inject({}) do |memo, object|
  #     memo[object.response_id] ||= []
  #     memo[object.response_id] << [object.text]
  #     memo
  #   end.each_pair.inject('') { |memo, obj| memo << { ResponseID: obj[0].to_s, Sentences: obj[1] }.to_json }
  # end

end

wt = Tree.new([], [], 'love')
puts "response_ids: " + wt.response_ids.to_s
puts "postfix_sets: " + wt.postfix_sets.to_s
