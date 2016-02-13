using System;
using System.Runtime.Serialization;

namespace TriggerManager.Repository
{
    [Serializable]
    public class RepositoryException : Exception
    {
        public RepositoryException() { }

        public RepositoryException(string message) : base(message) { }

        public RepositoryException(string message, Exception innerException) : base(message, innerException) { }

        protected RepositoryException(SerializationInfo serializationInfo, StreamingContext streamingContext) : base(serializationInfo, streamingContext) { }
    }
}
