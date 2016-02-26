# -*- encoding : utf-8 -*-
require_relative './base'

class Tree < Base
  attr_accessor :target, :responses, :comments, :response_ids, :postfix_sets

  def initialize(responses, comments, target)
    self.responses  = responses
    self.comments   = comments
    self.target     = target

    parse_words
  end

  def parse_words
    self.postfix_sets = JSON.parse(golangPostfixSets(json_sentences, target))
    self.response_ids = postfix_sets.map{|item| item['ResponseID']}
  end

  private

  # #  PRODUCTION
  # def json_sentences
  #   (responses + comments).as_json.to_json
  # end

end

wt = Tree.new([], [], "love")
puts "response_ids: " + wt.response_ids.to_s
puts "postfix_sets: " + wt.postfix_sets.to_s
