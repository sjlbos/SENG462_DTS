
using System;

namespace RabbitMessaging
{
    public interface IMessageReceiver : IDisposable
    {
        string LastMessage { get; }

        string GetNextMessage();

        void AckLastMessage();
        void NackLastMessage();

        void CancelMessageWait();
    }
}
