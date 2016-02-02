
using System;
using System.Runtime.Serialization;

namespace TransactionEvents
{
    [Serializable]
    public class UnrecognizedTransactionEventException : Exception
    {
        public UnrecognizedTransactionEventException() { }

        public UnrecognizedTransactionEventException(string message) : base(message) { }

        public UnrecognizedTransactionEventException(string message, Exception innerException) : base(message, innerException) { }

        protected UnrecognizedTransactionEventException(SerializationInfo serializationInfo, StreamingContext streamingContext) : base(serializationInfo, streamingContext) { }
    }
}
