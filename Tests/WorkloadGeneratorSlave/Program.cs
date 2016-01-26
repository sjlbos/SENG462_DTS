
using System;
using System.Configuration;
using log4net;
using ServiceHost;

namespace WorkloadGeneratorSlave
{
    public class Program
    {
        static void Main()
        {
            ILog log = LogManager.GetLogger(typeof (Program));
            try
            {
                string serviceName = ConfigurationManager.AppSettings["ServiceName"];
                string displayName = ConfigurationManager.AppSettings["DisplayName"];
                string serviceDescription = ConfigurationManager.AppSettings["ServiceDescription"];
                ServiceRunner.RunService<WorkloadGeneratorService>(serviceName, displayName, serviceDescription);
            }
            catch (Exception ex)
            {
                log.Fatal("The Workload Generator slave encountered a fatal error.", ex); 
            }
        }
    }
}
