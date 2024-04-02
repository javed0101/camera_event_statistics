import requests
import time

def hit_curl(camera_id, event_type):
    url = f'http://localhost:8082/api/v1/camera/event/stats?cameraID={camera_id}&eventType={event_type}'
    response = requests.get(url)
    if response.status_code == 200:
        print(f"Request to {url} is successful with status code {response.status_code}.")
    else:
        print(f"Request to {url} failed with status code {response.status_code}.")
    resp_json = response.json()
    if 'event' in resp_json and 'requestID' in resp_json['event']:
        print(resp_json['event']['requestID'])

def hit_curl_at_rate(camera_id, event_type, rate_per_second, total_requests):
    interval = 1 / rate_per_second
    for _ in range(total_requests):
        hit_curl(camera_id, event_type)
        time.sleep(interval)

camera_id = "10094541s"
event_type = "motion"
rate_per_second = 100 
total_requests = 1000

hit_curl_at_rate(camera_id, event_type, rate_per_second, total_requests)
