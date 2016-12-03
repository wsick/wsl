using System;
using System.Collections.Generic;
using System.IO;
using System.Xml.Linq;
using System.Linq;

namespace Wsl.Encoding {
    internal class XmlToWslTranspiler {
        private string indentToken = "  ";

        public XNamespace DefaultNs;
        public string Indent;
        public TextWriter Writer;
        private Stack<XElement> resolvers = new Stack<XElement>();
        private XElement currentResolver;

        public void PushResolver(XElement resolver) {
            resolvers.Push(resolver);
            currentResolver = resolver;
        }
        public XElement PopResolver() {
            var last = resolvers.Pop();
            if (resolvers.Count >= 1)
                currentResolver = resolvers.Peek();
            return last;
        }
        public void WriteDocument(XDocument xd) {
            var els = xd.Elements().ToList();
            var many = els.Count > 0;
            foreach (var el in els) {
                WriteElement(el, many);
            }
        }
        public void WriteElement(XElement el, bool many) {
            PushResolver(el);
            try {
                Writer.Write(many ? Indent : " ");
                WriteName(el.Name);
                WriteBlock("(", ")", false,
                    el.Attributes().Where(attr => attr.IsNamespaceDeclaration).ToList(), 
                    (attr, amany) => WriteAttribute(attr.Name, attr.Value, amany));
                WriteBlock("[", "]", false,
                    el.Attributes().Where(attr => !attr.IsNamespaceDeclaration).ToList(), 
                    (attr, amany) => WriteAttribute(attr.Name, attr.Value, amany));
                WriteBlock("{", "}", true,
                    el.Elements().ToList(), 
                    (child, emany) => WriteElement(child, emany));
                if (many)
                    Writer.Write(Writer.NewLine);
            } finally {
                PopResolver();
            }
        }
        public void WriteDeclaration(XName name, string value, bool many) {
            Writer.Write(many ? Indent : " ");
            Writer.Write("wslns");
            if (name.LocalName != "xmlns") {
                Writer.Write(":");
                Writer.Write(name.LocalName);
            }
            Writer.Write("=\"");
            Writer.Write(value);
            Writer.Write("\"");
            if (many)
                Writer.Write(Writer.NewLine);
        }
        public void WriteAttribute(XName name, string value, bool many) {
            Writer.Write(many ? Indent : " ");
            WriteName(name);
            Writer.Write("=\"");
            Writer.Write(value);
            Writer.Write("\"");
            if (many)
                Writer.Write(Writer.NewLine);
        }
        public void WriteName(XName name) {
            if (name.NamespaceName != DefaultNs.NamespaceName) {
                var prefix = currentResolver != null ? currentResolver.GetPrefixOfNamespace(name.NamespaceName) : "";
                if (!string.IsNullOrEmpty(prefix)) {
                    Writer.Write(prefix);
                    Writer.Write(":");
                }
            }
            WriteWslLocalName(name.LocalName);
        }
        private void WriteWslLocalName(string name) {
            var prev = ' ';
            foreach (var c in name) {
                if (char.IsUpper(c)) {
                    switch (prev) {
                        case ' ':
                        case '.':
                        case '-':
                            break;
                        default:
                            Writer.Write("-");
                            break;
                    }
                    Writer.Write(char.ToLower(c));
                } else {
                    Writer.Write(c);
                }
                prev = c;
            }
        }
        public void IncreaseIndent() {
            Indent += indentToken;
        }
        public void DecreaseIndent() {
            if (Indent.Length >= indentToken.Length) {
                Indent = Indent.Substring(0, Indent.Length - indentToken.Length);
            }
        }

        private void WriteBlock<T>(string begin, string end, bool forceMany, List<T> items, Action<T, bool> doEach) {
            var count = items.Count;
            if (count < 1) 
                return;
            var many = count > 1 || forceMany;
            Writer.Write(" ");
            Writer.Write(begin);
            IncreaseIndent();
            if (many)
                Writer.Write(Writer.NewLine);
            items.ForEach(item => doEach(item, many));
            DecreaseIndent();
            if (!many)
                Writer.Write(" ");
            else
                Writer.Write(Indent);
            Writer.Write(end);
        }
    }
}