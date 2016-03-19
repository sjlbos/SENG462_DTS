
using System;
using System.Runtime.Serialization;

namespace TriggerManager
{
    [Serializable]
    public class TriggerCommitException : Exception
    {
        public TriggerCommitException() { }

        public TriggerCommitException(string message) : base(message) { }

        public TriggerCommitException(string message, Exception innerException) : base(message, innerException) { }

        protected TriggerCommitException(SerializationInfo serializationInfo, StreamingContext streamingContext) : base(serializationInfo, streamingContext) { }
    }
}
