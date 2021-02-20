// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.


using System;
using System.Text;
using System.Security.Cryptography;
using Dolittle.SDK.Events;

namespace Dolittle.Data.Backups.Domain
{
    public static class EventSources
    {
        public static EventSourceId From(Guid application, string environment)
            => MergeTwoGuids(application, GuidFromString(environment));

        static Guid GuidFromString(string input)
        {
            using var md5 = MD5.Create();
            var hash = md5.ComputeHash(Encoding.Default.GetBytes(input));
            return new Guid(hash);
        }
        /// https://stackoverflow.com/a/1641173
        static Guid MergeTwoGuids(Guid guid1, Guid guid2)
        {
            const int BYTECOUNT = 16;
            var destByte = new byte[BYTECOUNT];
            var guid1Byte = guid1.ToByteArray();
            var guid2Byte = guid2.ToByteArray();
            
            for (var i = 0; i < BYTECOUNT; i++)
            {
                destByte[i] = (byte) (guid1Byte[i] ^ guid2Byte[i]);
            }
            return new Guid(destByte);
        }
    }
}
