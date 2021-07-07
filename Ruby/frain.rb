require 'uri'
require 'json'
require 'net/http'
require 'active_support'
require 'active_support/core_ext/numeric/time'

class Frain 
  class FrainError < StandardError
  end
  
  class ConnectionError < FrainError
  end
  
  class ServerError < FrainError
  end
  
  ENDPOINT = 'https://api.frain.dev/api/v1/banks'
  
  def initialize(redis)
    @redis_client = redis
    @expiry_time = 1.minute
  end
  
  def get_banks
    cached_data = fetch_from_cache
    
    response = case
               when cached_data.nil?
                 response = fetch_from_api
                 save_to_cache(response)
               else
                 cached_data
               end
    response
  end
  
  private
  
  def fetch_from_api
    uri = URI(ENDPOINT)
    request = Net::HTTP::Get.new(uri.request_uri)
    request['Authorization'] = ENV['FRAIN_API_TOKEN']
    
    http_client = Net::HTTP.new(uri.host, uri.port)
    http_client.use_ssl = true
    http_client.open_timeout = 10
    http_client.continue_timeout = 10
    http_client.read_timeout = 10
    http_client.ssl_timeout = 10
    
    begin 
      response = http_client.request(request)
    rescue => e
      raise ConnectionError.new("Could not connect to Frain, with error '#{e.message}'.")
    end
    
    result = {}
    if response.code.to_i == 200
      begin
        result = JSON.parse(response.body)
      rescue JSON::JSONError
        raise ServerError.new("Could not interpret Frain server response: '#{response.body}'.")
      end
    end
    
    if result['status'] != true
      raise ServerError.new("Could not read from Frain, server responded with '#{response.code}' returning: '#{response.body}'")
    end
    
    result['data']
  end
  
  def fetch_from_cache
    cached_data = @redis_client.get(banks_key)
    cached_data.nil? ? nil : JSON.parse(@redis_client.get(banks_key))
  end
  
  def save_to_cache(response)
    return response if response.empty?
    @redis_client.set(banks_key, JSON.dump(response), @expiry_time)
    response
  end
  
  def banks_key
    :frain_banks_uptime_status
  end
end

# TODO: Remove this test.
class RedisCache
  def initialize
    @data = {}
  end

  def get(key)
    @data[key]
  end

  def set(key, data, expiry)
    @data[key] = data
  end
end

redis = RedisCache.new
frain = Frain.new(redis)
pp frain.get_banks
