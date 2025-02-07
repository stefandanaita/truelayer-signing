package tlsigning

import (
	"io/ioutil"
	"testing"

	"github.com/Truelayer/truelayer-signing/go/errors"
	"github.com/stretchr/testify/assert"
)

const (
	Kid = "45fc75cf-5649-4134-84b3-192c2c78e990"
)

func getTestKeys(assert *assert.Assertions) ([]byte, []byte) {
	privateKeyBytes, err := ioutil.ReadFile("./testdata/ec512-private.pem")
	assert.Nilf(err, "private key read failed: %v", err)
	publicKeyBytes, err := ioutil.ReadFile("./testdata/ec512-public.pem")
	assert.Nilf(err, "public key read failed: %v", err)
	return privateKeyBytes, publicKeyBytes
}

func TestVerifyV1StaticSignatureShouldFail(t *testing.T) {
	assert := assert.New(t)

	_, publicKeyBytes := getTestKeys(assert)

	tlSignature := "eyJhbGciOiJFUzUxMiIsImtpZCI6IjQ1ZmM3NWNmLTU2NDktNDEzNC04NGIzLTE5MmMyYzc4ZTk5MCJ9..AdDESSiHQVQSRFrD8QO6V8m0CWIfsDGyMOlipOt9LQhyG1lKjDR17crBgy_7TYi4ZQH--dyNtN9Nab3P7yFQzgqOALl8S-beevWYpnIMXHQCgrv-XpfNtenJTckCH2UAQIwR-pjV8XiTM1be1RMYpMl8qYTbCL5Bf8t_dME-1E6yZQEH"
	body := []byte("{\"abc\":123}")

	err := VerifyWithPem(publicKeyBytes).
		Body(body).
		Verify(tlSignature)
	assert.NotNilf(err, "v1 signature verification should fail: %v", err)
	assert.ErrorAs(&errors.JwsError{}, &err, "error should be a JwsError")
}

func TestSignature(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("POST").
		Path(path).
		RequireHeader("Idempotency-Key").
		Header("X-Whatever-2", []byte("t2345d")).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Verify(signature)
	assert.Nilf(err, "signature verification should not fail: %v", err)
}

func TestVerifyStaticSignature(t *testing.T) {
	assert := assert.New(t)

	_, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"
	tlSignature := "eyJhbGciOiJFUzUxMiIsImtpZCI6IjQ1ZmM3NWNmLTU2NDktNDEzNC04NGIzLTE5MmMyYzc4ZTk5MCIsInRsX3ZlcnNpb24iOiIyIiwidGxfaGVhZGVycyI6IklkZW1wb3RlbmN5LUtleSJ9..AfhpFccUCUKEmotnztM28SUYgMnzPNfDhbxXUSc-NByYc1g-rxMN6HS5g5ehiN5yOwb0WnXPXjTCuZIVqRvXIJ9WAPr0P9R68ro2rsHs5HG7IrSufePXvms75f6kfaeIfYKjQTuWAAfGPAeAQ52PNQSd5AZxkiFuCMDvsrnF5r0UQsGi"

	err := VerifyWithPem(publicKeyBytes).
		Method("POST").
		Path(path).
		Header("X-Whatever-2", []byte("t2345d")).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Verify(tlSignature)
	assert.Nilf(err, "signature verification should not fail: %v", err)
}

func TestSignatureMethodMismatch(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("DELETE"). // different
		Path(path).
		Header("X-Whatever-2", []byte("aoitbeh")).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Verify(signature)
	assert.NotNilf(err, "signature verification should fail: %v", err)
	assert.ErrorAs(&errors.JwsError{}, &err, "error should be a JwsError")
}

func TestSignatureHeaderMismatch(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("post").
		Path(path).
		Header("X-Whatever-2", []byte("aoitbeh")).
		Header("Idempotency-Key", []byte("something-else")).
		Body(body).
		Verify(signature)
	assert.NotNilf(err, "signature verification should fail: %v", err)
	assert.ErrorAs(&errors.JwsError{}, &err, "error should be a JwsError")
}

func TestSignatureBodyMismatch(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("post").
		Path(path).
		Header("X-Whatever-2", []byte("aoitbeh")).
		Header("Idempotency-Key", idempotencyKey).
		Body([]byte("{\"max_amount_in_minor\":1234}")). // different
		Verify(signature)
	assert.NotNilf(err, "signature verification should fail: %v", err)
	assert.ErrorAs(&errors.JwsError{}, &err, "error should be a JwsError")
}

func TestSignatureMissingSignatureHeader(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("post").
		Path(path).
		Header("X-Whatever-2", []byte("aoitbeh")).
		// missing Idempotency-Key
		Body(body).
		Verify(signature)
	assert.NotNilf(err, "signature verification should fail: %v", err)
	assert.ErrorAs(&errors.JwsError{}, &err, "error should be a JwsError")
}

func TestRequiredHeaderMissingFromSignature(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("post").
		Path(path).
		RequireHeader("X-Required").
		Header("Idempotency-Key", idempotencyKey).
		Body(body).
		Verify(signature)
	assert.NotNilf(err, "signature verification should fail: %v", err)
	assert.ErrorAs(&errors.InvalidKeyError{}, &err, "error should be an InvalidKeyError")
}

func TestFlexibleHeaderCaseOrderVerify(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	idempotencyKey := []byte("idemp-2076717c-9005-4811-a321-9e0787fa0382")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Header("Idempotency-Key", idempotencyKey).
		Header("X-Custom", []byte("123")).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("POST").
		Path(path).
		Header("X-CUSTOM", []byte("123")).
		Header("idempotency-key", idempotencyKey).
		Body(body).
		Verify(signature)
	assert.Nilf(err, "signature verification should not fail: %v", err)
}

func TestEnforceDetached(t *testing.T) {
	assert := assert.New(t)

	_, publicKeyBytes := getTestKeys(assert)

	// signature for `/bar` but with a valid jws-body pre-attached
	tlSignature := "eyJhbGciOiJFUzUxMiIsImtpZCI6IjQ1ZmM3NWNmLTU2NDktNDEzNC04NGIzLTE5MmMyYzc4ZTk5MCIsInRsX3ZlcnNpb24iOiIyIiwidGxfaGVhZGVycyI6IiJ9.UE9TVCAvYmFyCnt9.ARLa7Q5b8k5CIhfy1qrS-IkNqCDeE-VFRDz7Lb0fXUMOi_Ktck-R7BHDMXFDzbI5TyaxIo5TGHZV_cs0fg96dlSxAERp3UaN2oCQHIE5gQ4m5uU3ee69XfwwU_RpEIMFypycxwq1HOf4LzTLXqP_CDT8DdyX8oTwYdUBd2d3D17Wd9UA"

	body := []byte("{}")
	path := "/foo"
	err := VerifyWithPem(publicKeyBytes).
		Method("post").
		Path(path).
		Body(body).
		Verify(tlSignature)
	assert.NotNilf(err, "signature verification should fail: %v", err)
	assert.ErrorAs(&errors.JwsError{}, &err, "error should be a JwsError")
}

func TestJwsHeaderExtraction(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, _ := getTestKeys(assert)

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("delete").
		Path("/foo").
		Header("X-Custom", []byte("123")).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	jwsHeader, _ := ExtractJwsHeader(signature)

	assert.Equal(jwsHeader.Alg, "ES512")
	assert.Equal(jwsHeader.Kid, Kid)
	assert.Equal(jwsHeader.TlVersion, "2")
	assert.Equal(jwsHeader.TlHeaders, "X-Custom")
}

func TestSignatureNoHeaders(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Method("post").
		Path(path).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Method("post").
		Path(path).
		Header("X-Whatever", []byte("aoitbeh")).
		Body(body).
		Verify(signature)
	assert.Nilf(err, "signature verification should not fail: %v", err)
}

func TestVerifyWithoutMethodShouldFail(t *testing.T) {
	assert := assert.New(t)

	privateKeyBytes, publicKeyBytes := getTestKeys(assert)

	body := []byte("{\"currency\":\"GBP\",\"max_amount_in_minor\":5000000}")
	path := "/merchant_accounts/a61acaef-ee05-4077-92f3-25543a11bd8d/sweeping"

	signature, err := SignWithPem(Kid, privateKeyBytes).
		Path(path).
		Body(body).
		Sign()
	assert.Nilf(err, "signing failed: %v", err)

	err = VerifyWithPem(publicKeyBytes).
		Path(path).
		Body(body).
		Verify(signature)
	assert.NotNilf(err, "signature verification should fail: %v", err)
	assert.ErrorAs(&errors.JwsError{}, &err, "error should be a JwsError")
}
