
using System;

namespace RabbitMessaging
{
    public interface IMessagePublisher : IDisposable
    {
        void PublishMessage(string message, string routingKey);
    }
}
