import argparse
import queue
import pika


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
	parts = line.rstrip().split(",")
	#print(parts)
	command = parts[0]
	userId = parts[1]

#Store in queues
	if userId not in UserList and command != 'DUMPLOG':
		UserList[userId] = list()
	elif command != 'DUMPLOG':
		UserList[userId].append(line.rstrip())
for user in UserList:
	for command in UserList[user]:
		channel.basic_publish(exchange='', routing_key='UserInputs', body=command)
		print("[x] Sent: " + command)






