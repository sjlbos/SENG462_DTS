using System;
using System.Runtime.Serialization;

namespace TriggerManager.Repository
{
    [Serializable]
    public class QuoteException : Exception
    {
        public QuoteException() { }

        public QuoteException(string message) : base(message) { }

        public QuoteException(string message, Exception innerException) : base(message, innerException) { }

        protected QuoteException(SerializationInfo serializationInfo, StreamingContext streamingContext) : base(serializationInfo, streamingContext) { }
    }
}
