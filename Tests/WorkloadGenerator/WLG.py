import argparse
import queue
import pika
import json


parser = argparse.ArgumentParser(description='Workload Generator for Distributed System')
parser.add_argument('filename', nargs='?')

connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost', port=5672))
channel = connection.channel()

channel.queue_declare(queue='UserInputs')

args = parser.parse_args()
filename = ""
if args.filename:
	filename = args.filename
else:
	filename = input("Workload File: ")


#Open File
fp = open(filename,'r')
UserList = dict()
for line in fp:
	line_parts = line.split(" ")
	command_line = line_parts[1]
	print(command_line)
	command_parts = command_line.rstrip().split(",")
	command= command_parts[0]
	userId = command_parts[1]
#Store in queues
	if userId not in UserList and command != 'DUMPLOG':
		UserList[userId] = queue.Queue()
	if command != 'DUMPLOG':
		UserList[userId].put(command_line)

sent_messages=0
for userId in UserList:
	UserCommands = list()
	userQueue = UserList[userId]
	while not userQueue.empty():
		command = userQueue.get()
		UserCommands.append(command)
	json_send = json.dumps(UserCommands, ensure_ascii=False)
	channel.basic_publish(exchange='', routing_key='UserInputs', body=json_send)
	sent_messages = sent_messages + 1
print("[x] Sent " + str(sent_messages) + " Messages")






