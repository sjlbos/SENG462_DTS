import pika

connection = pika.BlockingConnection(pika.ConnectionParameters(
        host='localhost', port=5672))
channel = connection.channel()

channel.queue_declare(queue='UserInputs')

def callback(ch, method, properties, body):
    print(" [x] Received %r" % body)

channel.basic_consume(callback,
                      queue='UserInputs',
                      no_ack=True)

print(' [*] Waiting for messages. To exit press CTRL+C')
channel.start_consuming()