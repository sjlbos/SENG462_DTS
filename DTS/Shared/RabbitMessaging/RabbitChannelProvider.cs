
using System;
using System.Collections.Concurrent;
using RabbitMQ.Client;

namespace RabbitMessaging
{
    internal class RabbitChannelProvider
    {
        private static readonly ConcurrentDictionary<HostConfiguration, IConnection> ConnectionPool = new ConcurrentDictionary<HostConfiguration, IConnection>(); 

        public static IModel OpenChannelToHost(HostConfiguration hostConfig)
        {
            if(hostConfig == null)
                throw new ArgumentNullException("hostConfig");

            if (ConnectionPool.ContainsKey(hostConfig) && ConnectionPool[hostConfig].IsOpen)
            {
                return ConnectionPool[hostConfig].CreateModel();
            }

            var newConnection = OpenNewConnection(hostConfig);
            ConnectionPool.AddOrUpdate(hostConfig, newConnection, (config, oldValue) => newConnection);
            var channel = newConnection.CreateModel();
            newConnection.AutoClose = true; // Set the connection to auto close after opening first channel
            return channel;
        }

        private static IConnection OpenNewConnection(HostConfiguration hostConfig)
        {
            if (hostConfig == null)
                throw new ArgumentNullException("hostConfig");

            var factory = new ConnectionFactory
            {
                HostName = hostConfig.HostName,
                VirtualHost = hostConfig.VirtualHost,
                Port = hostConfig.Port,
                UserName = hostConfig.Username,
                Password = hostConfig.Password
            };

            return factory.CreateConnection();
        }
    }
}
