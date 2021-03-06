import argparse
import subprocess
import pika
import json
import Queue
import urllib2
import random
import sys


parser = argparse.ArgumentParser(description='Workload Generator for Distributed System')
parser.add_argument('--filename', nargs='?')
parser.add_argument('--hostname', nargs='?')
parser.add_argument('--port', nargs='?')
parser.add_argument('--slaves', nargs='?')
parser.add_argument('--rhost', nargs='?')
parser.add_argument('--rport', nargs='?')
parser.add_argument('--nhost', nargs='?')
parser.add_argument('--nport', nargs='?')

args = parser.parse_args()

rabbitHost = "localhost"
rabbitPort = "44410"
nannyHost = "localhost"
nannyPort = "44413"
filename = ""
outputFile = ""
hostname = "localhost"
num_Slaves = 1
doDump = False
port = ""

workerSum=0
commandSum=0

if args.filename:
	filename = args.filename
else:
	filename = raw_input("Workload File: ")
#if args.hostname:
	#hostname = args.hostname
if args.port:
	port = args.port
if args.slaves:
	num_Slaves = args.slaves
if args.rhost:
	rabbitHost = args.rhost
if args.rport:
	rabbitPort = args.rport
if args.nhost:
	nannyHost = args.nhost
if args.rport:
	nannyPort = args.nport

hostname = ["b137.seng.uvic.ca", "b153.seng.uvic.ca"]

urls = dict()
urls[0] = "http://" + hostname[0]
urls[1] = "http://" + hostname[1]
if port:
	urls[0] = url[0] + ":" + port
	urls[1] = url[1] + ":" + port

credentials = pika.PlainCredentials('dts_user', 'Group1')
connection = pika.BlockingConnection(pika.ConnectionParameters(rabbitHost, int(rabbitPort), '/', credentials))
channel = connection.channel()

channel.exchange_declare(exchange='WorkloadGenerator',type='direct', durable=True, auto_delete=True)
for i in range(1, int(num_Slaves)+1):
	channel.queue_declare(queue='Slave' +str(i), durable=True)
channel.queue_declare(queue='SlaveStatus', durable=True)
channel.queue_bind(exchange='WorkloadGenerator',
                   queue='SlaveStatus',
                   routing_key='SlaveStatus')
class ApiCommand:
	def __init__(self, uri, request, newId, method, statuscode):
		self.Uri = uri
		self.RequestBody = request
		self.Id = newId
		self.Method = method
		self.ExpectedStatusCode = statuscode

	def reprJSON(self):
		return dict(Uri=self.Uri, RequestBody=self.RequestBody, Id=self.Id, Method=self.Method, ExpectedStatusCode=self.ExpectedStatusCode)

class BatchCommand:
	def __init__(self,Type, Id):
		self.MessageType = Type
		self.Id = Id
		self.Commands = []

	def add_command(self, Command):
		self.Commands.append(Command)

	def reprJSON(self):
		return dict(MessageType=self.MessageType, Id=self.Id, Commands=self.Commands)

class ControlCommand:
	def __init__(self):
		self.Id = "Start"
		self.MessageType = "Control"
		self.Command = "Start"

	def reprJSON(self):
		return dict(Id=self.Id, MessageType=self.MessageType, Command=self.Command)

class ComplexEncoder(json.JSONEncoder):
    def default(self, obj):
        if hasattr(obj,'reprJSON'):
            return obj.reprJSON()
        else:
            return json.JSONEncoder.default(self, obj)

def encode_Command(obj):
    if isinstance(obj, Command):
        return obj.__dict__
    return obj

def getAddCommand(User, Amount, Id):
	uri = urls[random.randint(0,1)] + "/api/users/" + User
	json_string = '{ "Amount" : "' + Amount + '" }'
	return ApiCommand(uri, json_string, Id, "PUT", 200)

def getQuoteCommand(User, StockSymbol, Id):
	uri = urls[random.randint(0,1)] + "/api/users/" + User + "/stocks/quote/" + StockSymbol
	return ApiCommand(uri, "", Id, "GET", 200)

def getBuyCommand(User, StockSymbol, Amount, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/pending-purchases"
	json_string = '{"Symbol" : "' + StockSymbol + '", "Amount" : "' + Amount + '"" }'
	return ApiCommand(uri, json_string, Id, "POST", 200)

def getCommitBuyCommand(User, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/pending-purchases/commit"
	return ApiCommand(uri, "", Id, "POST", 200)

def getCancelBuyCommand(User, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/pending-purchases"
	return ApiCommand(uri, "", Id, "DELETE", 200)

def getSellCommand(User, StockSymbol, Amount, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/pending-sales"
	json_string = '{"Symbol" : "' + StockSymbol + '", "Amount" : "' + Amount + '" }'
	return ApiCommand(uri, json_string, Id, "POST", 200)

def getCommitSellCommand(User, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/pending-sales/commit"
	return ApiCommand(uri, "", Id, "POST", 200)

def getCancelSellCommand(User, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/pending-sales"
	return ApiCommand(uri, "", Id, "DELETE", 200)

def getSetBuyAmountCommand(User, StockSymbol, Amount, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/buy-triggers/"+StockSymbol
	json_string = '{"Amount" : "' + Amount + '" }'
	return ApiCommand(uri, json_string, Id, "PUT", 200)

def getCancelSetBuyCommand(User, StockSymbol, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/buy-triggers/"+StockSymbol
	return ApiCommand(uri, "", Id, "DELETE", 200)

def getSetBuyTriggerCommand(User, StockSymbol, Price, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/buy-triggers/"+StockSymbol
	json_string = '{"Price" : "' + Price + '" }'
	return ApiCommand(uri, json_string, Id, "PUT", 200)

def getSetSellAmountCommand(User, StockSymbol, Amount, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/sell-triggers/"+StockSymbol
	json_string = '{"Amount" : "' + Amount + '" }'
	return ApiCommand(uri, json_string, Id, "PUT", 200)

def getSetSellTriggerCommand(User, StockSymbol, Price, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/sell-triggers/"+StockSymbol
	json_string = '{"Price" : "' + Price + '"}'
	return ApiCommand(uri, json_string, Id, "PUT", 200)

def getCancelSetSellCommand(User, StockSymbol, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/sell-triggers/"+StockSymbol
	return ApiCommand(uri, "", Id, "DELETE", 200)

def getDumplogUserCommand(User,Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/transactions"
	return ApiCommand(uri, "", Id, "GET", 200)

def getDumplogCommand(Id):
	uri = urls[random.randint(0,1)] + "/api/users/transactions"
	return ApiCommand(uri, "", Id, "GET", 200)

def getDisplaySummaryCommand(User, Id):
	uri = urls[random.randint(0,1)] + "/api/users/"+User+"/summary"
	return ApiCommand(uri, "", Id, "GET", 200)

#Open File
fp = open(filename,'r')
UserList = dict()
for line in fp:
	commandSum = commandSum + 1
	line_parts = line.split(" ")
	CommandNo = line_parts[0]
	command_line = line_parts[1]
	#print(command_line)
	command_parts = command_line.rstrip().split(",")
	command= command_parts[0]
	userId = command_parts[1]

#Store in queues
	if userId not in UserList and command != 'DUMPLOG':
		UserList[userId] = Queue.Queue()
	if command != 'DUMPLOG':
		UserList[userId].put(line)

	if command == "DUMPLOG":
		if len(command_parts) == 2:
			doDump = True
		else:
			outputFile = command_parts[2]
			if userId not in UserList:
				UserList[userId] = queue.Queue()
			UserList[userId].put(line)



sent_messages=0
for userId in UserList:
	UserCommands = BatchCommand("BatchOrder",userId)
	userQueue = UserList[userId]
	while not userQueue.empty():
		line = userQueue.get()
		line_parts = line.split(" ")
		TransId = line_parts[0].lstrip('[').rstrip(']')
		command = line_parts[1].split(",")
		tmpCommand = None;
		if command[0] == "ADD":
			tmpCommand = getAddCommand(command[1], command[2], TransId)
		elif command[0] == "QUOTE":
			tmpCommand = getQuoteCommand(command[1], command[2], TransId)
		elif command[0] == "BUY":
			tmpCommand = getBuyCommand(command[1], command[2], command[3], TransId)
		elif command[0] == "COMMIT_BUY":
			tmpCommand = getCommitBuyCommand(command[1], TransId)
		elif command[0] == "CANCEL_BUY":
			tmpCommand = getCancelBuyCommand(command[1], TransId)
		elif command[0] == "SELL":
			tmpCommand = getSellCommand(command[1], command[2], command[3], TransId)
		elif command[0] == "COMMIT_SELL":
			tmpCommand = getCommitSellCommand(command[1], TransId)
		elif command[0] == "CANCEL_SELL":
			tmpCommand = getCancelSellCommand(command[1], TransId)
		elif command[0] == "SET_BUY_AMOUNT":
			tmpCommand = getSetBuyAmountCommand(command[1], command[2], command[3], TransId)
		elif command[0] == "CANCEL_SET_BUY":
			tmpCommand = getCancelSetBuyCommand(command[1], command[2], TransId)
		elif command[0] == "SET_BUY_TRIGGER":
			tmpCommand = getSetBuyTriggerCommand(command[1], command[2], command[3], TransId)
		elif command[0] == "SET_SELL_AMOUNT":
			tmpCommand = getSetSellAmountCommand(command[1], command[2], command[3], TransId)
		elif command[0] == "CANCEL_SET_SELL":
			tmpCommand = getCancelSetSellCommand(command[1], command[2], TransId)
		elif command[0] == "SET_SELL_TRIGGER":
			tmpCommand = getSetSellTriggerCommand(command[1], command[2], command[3], TransId)
		elif command[0] == "DISPLAY_SUMMARY":
			tmpCommand = getDisplaySummaryCommand(command[1], TransId)
		elif command[0] == "DUMPLOG":
			tmpCommand = getDumplogUserCommand(command[1], TransId)
		else:
			UserCommands.add_command(getAddCommand("","",0))
		UserCommands.add_command(tmpCommand)
	json_send = json.dumps(UserCommands.reprJSON(), cls=ComplexEncoder)
	
	slaveNo = sent_messages % int(num_Slaves) + 1
	rKey = "Slave" + str(slaveNo)
	rExchange = "WorkloadGenerator"
	print("Sending to Slave " + rKey)
	channel.basic_publish(exchange=rExchange, routing_key=rKey, body=json_send)
	sent_messages = sent_messages + 1

print("Sent " + str(commandSum) + " Commands")
raw_input("Press Enter to continue...")
messageCommand = ControlCommand()
json_send = json.dumps(messageCommand.reprJSON(), cls=ComplexEncoder)
print(json_send)
channel.basic_publish(exchange='WorkloadGenerator', routing_key='Control', body=json_send)

def callback(ch, method, properties, body):
    print(" [x] Received %r" % body)
    global workerSum
    global sent_messages
    global num_Slaves
    global outputFile
    jsonVal = json.loads(body)
    print(" [x] Done")
    print jsonVal
    workerNum = jsonVal['SlaveName']
    workerSum -= int(workerNum)
    ch.basic_ack(delivery_tag = method.delivery_tag)
    if workerSum == 0:
    	print("Complete")
    	if doDump:
		print("Getting XML File")
		subprocess.call("wget http://b136.seng.uvic.ca:44410/audit/transactions/testLog", shell=True)
		sys.exit(0)


print("Waiting for workers")
for i in range(int(num_Slaves)):
	workerSum += (i + 1)
print("current workerSum: " + str(workerSum))
channel.basic_qos(prefetch_count=1)
channel.basic_consume(callback, queue='SlaveStatus')

channel.start_consuming()

#{ "Status" : "Complete" , "WorkerNum" : int} 


