import argparse

parser = argparse.ArgumentParser(description='Workload Generator for Distributed System')
parser.add_argument('filename', nargs='?')

args = parser.parse_args()
filename = ""

if args.filename:
	filename = args.filename
else:
	filename = input("Workload File: ")

filenameOut = filename + 'Updated'

fp = open(filename,'r')
fout = open(filenameOut, 'w')

for line in fp:
	parts = line.split(' ')
	tmpString = parts[1]
	fout.write(tmpString + '\n')
