using Microsoft.Build.Utilities;
using Microsoft.Build.Framework;
using Microsoft.Web.XmlTransform;

namespace MSBuildTasks
{
    public class TransformXml : Task
    {
        [Required]
        public ITaskItem Source { get; set; }

        [Required]
        public ITaskItem Transform { get; set; }

        [Required]
        public ITaskItem Destination { get; set; }

        public override bool Execute()
        {
            var originalFileXml = new System.Xml.XmlDocument();
            originalFileXml.Load(Source.ItemSpec);

            using (var xmlTransform = new XmlTransformation(Transform.ItemSpec))
            {
                if (xmlTransform.Apply(originalFileXml) == false)
                    return false;

                originalFileXml.Save(Destination.ItemSpec);
            }

            return true;
        }
    }
}