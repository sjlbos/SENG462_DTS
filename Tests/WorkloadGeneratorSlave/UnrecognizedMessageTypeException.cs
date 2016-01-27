using System;
using System.Runtime.Serialization;

namespace WorkloadGeneratorSlave
{
    [Serializable]
    public class UnrecognizedMessageTypeException : Exception
    {
        public UnrecognizedMessageTypeException() { }

        public UnrecognizedMessageTypeException(string message) : base(message) { }

        public UnrecognizedMessageTypeException(string message, Exception innerException) : base(message, innerException) { }

        protected UnrecognizedMessageTypeException(SerializationInfo serializationInfo, StreamingContext streamingContext) : base(serializationInfo, streamingContext) { }
    }
}
