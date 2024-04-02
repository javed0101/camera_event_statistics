import json, os, time, logging, random, uuid
from datetime import datetime, timedelta
from dotenv import load_dotenv
from pulsar import Client, Producer

load_dotenv()

hostname = os.getenv("HOSTNAME")
topic = os.getenv("TOPIC")
eventRate = int(os.getenv("EVENT_RATE"))
eventType = json.loads(os.getenv("EVENT_TYPE"))
cameraIDs = json.loads(os.getenv("CAMERA_IDS"))

directory = os.path.dirname(logFile)
os.makedirs(directory, exist_ok=True)

def push_dummy_camera_events(producer):
    for index in range(eventRate):
        timestamp = get_current_time()
        camera_event = {
            "eventID": str(uuid.uuid4()),
            "cameraID": random.choice(cameraIDs),
            "timestamp": timestamp,
            "location": {
                "latitude": random.uniform(-90.0, 90.0),
                "longitude": random.uniform(-180.0, 180.0)
            },
            "eventType": random.choice(eventType),
            "metaData": {
                "objectID": str(uuid.uuid4()),
                "confidenceScore": random.uniform(0, 1),
            } 
        }
        producer.send(json.dumps(camera_event).encode())
        print(camera_event)
        logging.info(f"Pushed event to Pulsar with event ID: {camera_event['eventID']} and event type: {camera_event['eventType']} at {time.time()}")

def get_current_time():
    current_time = datetime.now()
    return current_time.strftime('%Y-%m-%dT%H:%M:%SZ')

client = Client(hostname)
producer = client.create_producer(topic)

try:
    while True:
        push_dummy_camera_events(producer)
        time.sleep(1)
except KeyboardInterrupt:
    print("Script interrupted")
except Exception as e:
    logging.error(f"An error occurred: {str(e)}")

producer.close()
client.close()
