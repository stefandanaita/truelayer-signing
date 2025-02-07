using Xunit;
using TrueLayer.Signing;
using System;
using FluentAssertions;

namespace Tests
{
    public class UsageTest
    {
        // Note: Should be the same keys used in all lang tests
        //       so static signature tests ensure cross-lang consistency.
        internal const string PublicKey = "-----BEGIN PUBLIC KEY-----\n"
            + "MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBVIVnghUzHmCEZ3HNjDmaZMJ7UwZf\n"
            + "av2SYcEtbDQc4uPhiEwWoYZMxzgvsz1vVGkusfTIjcXeCfDZ+xu9grRYt4kBo39z\n"
            + "w0i0j1rau4T7Bi+thc/VZpCyuwt63mZWcRs5PlQzpL34bBSXL5L6G9XUtXn8pXwU\n"
            + "GMhNDp5xVGbslRqTU8s=\n"
            + "-----END PUBLIC KEY-----\n";
        internal const string PrivateKey = "-----BEGIN EC PRIVATE KEY-----\n"
            + "MIHcAgEBBEIAVItA/A9H8WA0rOmDO5kq774be6noZ73xWJkbmzihkhtnYJ+eCQl4\n"
            + "G68ZFKildLuR2DElMBrNgJHY1TkL9hr7U9GgBwYFK4EEACOhgYkDgYYABAFUhWeC\n"
            + "FTMeYIRncc2MOZpkwntTBl9q/ZJhwS1sNBzi4+GITBahhkzHOC+zPW9UaS6x9MiN\n"
            + "xd4J8Nn7G72CtFi3iQGjf3PDSLSPWtq7hPsGL62Fz9VmkLK7C3reZlZxGzk+VDOk\n"
            + "vfhsFJcvkvob1dS1efylfBQYyE0OnnFUZuyVGpNTyw==\n"
            + "-----END EC PRIVATE KEY-----\n";
        internal const string Kid = "45fc75cf-5649-4134-84b3-192c2c78e990";

        [Fact]
        public void SignAndVerify()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Sign();

            Verifier.VerifyWithPem(PublicKey)
                .Method("post") // case-insensitive: no troubles
                .Path(path)
                .Header("X-Whatever-2", "t2345d")
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Verify(tlSignature); // should not throw
        }

        [Fact]
        public void SignAndVerify_NoHeaders()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Body(body)
                .Sign();

            Verifier.VerifyWithPem(PublicKey)
                .Method("POST") 
                .Path(path)
                .Body(body)
                .Verify(tlSignature); // should not throw
        }

        // Verify the a static signature used in all lang tests to ensure
        // cross-lang consistency and prevent regression.
        [Fact]
        public void VerifyStaticSignature()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";
            var tlSignature = "eyJhbGciOiJFUzUxMiIsImtpZCI6IjQ1ZmM3NWNmLTU2NDktNDEzNC04NGIzLTE5MmMyYzc4ZTk5MCIsInRsX3ZlcnNpb24iOiIyIiwidGxfaGVhZGVycyI6IklkZW1wb3RlbmN5LUtleSJ9..AfhpFccUCUKEmotnztM28SUYgMnzPNfDhbxXUSc-NByYc1g-rxMN6HS5g5ehiN5yOwb0WnXPXjTCuZIVqRvXIJ9WAPr0P9R68ro2rsHs5HG7IrSufePXvms75f6kfaeIfYKjQTuWAAfGPAeAQ52PNQSd5AZxkiFuCMDvsrnF5r0UQsGi";

            Verifier.VerifyWithPem(PublicKey)
                .Method("POST")
                .Path(path)
                .Header("X-Whatever-2", "t2345d")
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Verify(tlSignature); // should not throw
        }

        [Fact]
        public void SignAndVerify_MethodMismatch()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Sign();

            Action verify = () => Verifier.VerifyWithPem(PublicKey)
                .Method("DELETE") // different
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Verify(tlSignature);

            verify.Should().Throw<SignatureException>();
        }

        [Fact]
        public void SignAndVerify_PathMismatch()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Sign();

            Action verify = () => Verifier.VerifyWithPem(PublicKey)
                .Method("POST")
                .Path("/merchant_accounts/67b5b1cf-1d0c-45d4-a2ea-61bdc044327c/sweeping") // different
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Verify(tlSignature);

            verify.Should().Throw<SignatureException>();
        }

        [Fact]
        public void SignAndVerify_HeaderMismatch()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Sign();

            Action verify = () => Verifier.VerifyWithPem(PublicKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", "something-else") // different
                .Body(body)
                .Verify(tlSignature);

            verify.Should().Throw<SignatureException>();
        }

        [Fact]
        public void SignAndVerify_BodyMismatch()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Sign();

            Action verify = () => Verifier.VerifyWithPem(PublicKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000001}") // different
                .Verify(tlSignature);

            verify.Should().Throw<SignatureException>();
        }

        [Fact]
        public void SignAndVerify_MissingSignatureHeader()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Sign();

            Action verify = () => Verifier.VerifyWithPem(PublicKey)
                .Method("POST")
                .Path(path)
                // missing Idempotency-Key
                .Body(body)
                .Verify(tlSignature);

            verify.Should().Throw<SignatureException>();
        }

        [Fact]
        public void SignAndVerify_RequiredHeaderMissingFromSignature()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Sign();

            Action verify = () => Verifier.VerifyWithPem(PublicKey)
                .Method("POST")
                .Path(path)
                .RequireHeader("X-Required") // missing from signature
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Verify(tlSignature);

            verify.Should().Throw<SignatureException>();
        }

        [Fact]
        public void SignAndVerify_FlexibleHeaderCaseOrderVerify()
        {
            var body = "{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}";
            var idempotency_key = "idemp-2076717c-9005-4811-a321-9e0787fa0382";
            var path = "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping";

            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("POST")
                .Path(path)
                .Header("Idempotency-Key", idempotency_key)
                .Header("X-Custom", "123")
                .Body(body)
                .Sign();

            Verifier.VerifyWithPem(PublicKey)
                .Method("POST")
                .Path(path)
                .Header("X-CUSTOM", "123") // different order & case, it's ok!
                .Header("Idempotency-Key", idempotency_key)
                .Body(body)
                .Verify(tlSignature);
        }

        [Fact]
        public void Verifier_ExtractKid()
        {
            var tlSignature = Signer.SignWithPem(Kid, PrivateKey)
                .Method("delete")
                .Path("/foo")
                .Header("X-Custom", "123")
                .Sign();

            Verifier.ExtractKid(tlSignature).Should().Be(Kid);
        }
    }
}
