
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Globalization;
using System.Net;
using System.Threading.Tasks;
using log4net;
using Newtonsoft.Json;
using RabbitMessaging;
using ServiceHost;

namespace WorkloadGeneratorSlave
{
    public class WorkloadQueueMonitor : QueueMonitorWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (WorkloadQueueMonitor));

        private readonly int _httpWorkerCount ;
        private readonly ConcurrentQueue<WorkloadBatchMessage> _batchQueue; 

        public WorkloadQueueMonitor(string instanceId, IMessageReceiver messageReceiver, int httpWorkerCount) 
            : base(instanceId, messageReceiver)
        {
            _batchQueue = new ConcurrentQueue<WorkloadBatchMessage>();
            _httpWorkerCount = httpWorkerCount;
        }

        public override void ProcessMessage(string message)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            try
            {
                var deserializedMessage = JsonConvert.DeserializeObject<WorkloadGeneratorMessage>(message, new WorkloadGeneratorMessageConverter());
                if (deserializedMessage is WorkloadBatchMessage)
                {
                    HandleApiCommandMessage(deserializedMessage as WorkloadBatchMessage);
                    return;
                }
                if (deserializedMessage is ControlMessage)
                {
                    HandleControlMessage(deserializedMessage as ControlMessage);
                    return;
                }      
            }
            catch (UnrecognizedMessageTypeException ex)
            {
                Log.Warn("Received a message of unrecognized type.", ex);  
            }
        }

        private void HandleApiCommandMessage(WorkloadBatchMessage message)
        {
            Log.DebugFormat(CultureInfo.InvariantCulture,
                "Received command batch with Id={0}.", message.Id);

            _batchQueue.Enqueue(message);
        }

        private void HandleControlMessage(ControlMessage message)
        {
            Log.InfoFormat(CultureInfo.InvariantCulture,
                "Received control message \"{0}\", Id={1}.", message.Command, message.Id);

            switch (message.Command)
            {
                case ControlMessage.StartCommand:
                    ProcessWorkloadOrders(_httpWorkerCount);
                    break;
                default:
                    Log.Error("Unrecognized command: " + message.Command);
                    break;
            }
        }

        private void ProcessWorkloadOrders(int workerThreadCount)
        {
            Log.InfoFormat(CultureInfo.InvariantCulture,
                "Starting command batch execution with {0} threads.", workerThreadCount);
            var taskList = new List<Task>();
            for (int i = 0; i < workerThreadCount; i++)
            {
                taskList.Add(Task.Run(() => ProcessWorkloadOrders()));
            }
            Task.WaitAll(taskList.ToArray());

            Log.Info("All command batch threads have finished executing.");
        }

        private void ProcessWorkloadOrders()
        {
            WorkloadBatchMessage batch = null;
            while (_batchQueue.TryDequeue(out batch))
            {
                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Worker batch with");
                foreach (var command in batch.Commands)
                {
                    ExecuteApiCommand(command);
                }
            }   
        }

        private void ExecuteApiCommand(ApiCommand command)
        {
            Log.DebugFormat(CultureInfo.InvariantCulture, 
                "Executing api command with Id={0}: \"{1} {2}\" with request body \"{3}\".", 
                command.Id, command.Method, command.Uri, command.RequestBody);
            try
            {
                using (var response = (HttpWebResponse) command.HttpRequest.GetResponse())
                {
                    if (response.StatusCode != command.ExpectedStatusCode)
                    {
                        Log.WarnFormat(CultureInfo.InvariantCulture,
                            "RESPONSE CODE ASSERTION FAILURE - Expected: {0} Received: {1}", command.ExpectedStatusCode,
                            response.StatusCode);
                    }
                    else
                    {
                        Log.DebugFormat("Api command with Id={0} executed successfully.", command.Id);
                    }
                }
            }
            catch (WebException ex)
            {
                Log.Error("Encountered an error while executing api command with Id=" + command.Id, ex);
            }
        }
    }
}
