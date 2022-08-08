package auth

import (
	"context"
	"fmt"
	"github.com/okcredit/go-common/encoding/json"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	HeaderAuthorization = "Authorization"
	TokenTypeBearer     = "Bearer"
	TokenTypeAdmin      = "X-Admin"
)

type contextKey int

/////////////////////////////////////////////////////
// Merchant Auth
/////////////////////////////////////////////////////
const (
	tokenEncryptionKey            = "cfjmghkgchfs8yfihft4ur865vhgfgjh"
	keyMerchantID      contextKey = 1
)

func VerifyMerchantToken(accessToken string) (string, error) {
	token, err := jwt.ParseEncrypted(accessToken)
	if err != nil {
		return "", err
	}

	claims := jwt.Claims{}
	if err := token.Claims([]byte(tokenEncryptionKey), &claims); err != nil {
		return "", err
	}

	expected := jwt.Expected{
		Audience: jwt.Audience([]string{}),
		Time:     time.Now(),
	}
	if err := claims.Validate(expected); err != nil {
		bytes, _ := json.Marshal(claims)
		return "", fmt.Errorf("invalid claims; err = %v; claims = %s", err, bytes)
	}

	return claims.Subject, nil
}

func WithMerchantContext(r *http.Request, merchantID string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, keyMerchantID, merchantID)
	r = r.WithContext(ctx)
	return r
}

func GetMerchantID(r *http.Request) string {
	if merchantID, ok := r.Context().Value(keyMerchantID).(string); ok {
		return merchantID
	} else {
		return ""
	}
}

/////////////////////////////////////////////////////
// Admin Auth
/////////////////////////////////////////////////////
const (
	publicKey  = `{"kty":"RSA","n":"yCTWGCwHQKoCW9wJ5O5OJkEECa21ZAR5UU9_HgIUg1eowdPSrVEZkRNVAwxAJ5JFh1vUMHoP-o-WVn9hevYQcq_pHG4616tGGCn6rFTak7HukLFEEnKJButdxoHxTKGLET-kXIZXnQ8iUOU-sZGEkVezHhMin-rkOEUQB6-rJ8OV-V5m8qjkzrcpJf_9Y_JkQweyomJ21RDsQUUG2rgclkjTJF_PZzklffUS8-gc3Mka0zYnEhzhWFUh0pcxpDpjcFA0WWSdUWW_ASVJFXlbYf0Fos_1IVEVs7svfeTq0w8N2H0XXC4gOHIX7qf7heXaLTJFsuMbs3HNoTP0wrwUtQ","e":"AQAB"}`
	privateKey = `{"kty":"RSA","n":"yCTWGCwHQKoCW9wJ5O5OJkEECa21ZAR5UU9_HgIUg1eowdPSrVEZkRNVAwxAJ5JFh1vUMHoP-o-WVn9hevYQcq_pHG4616tGGCn6rFTak7HukLFEEnKJButdxoHxTKGLET-kXIZXnQ8iUOU-sZGEkVezHhMin-rkOEUQB6-rJ8OV-V5m8qjkzrcpJf_9Y_JkQweyomJ21RDsQUUG2rgclkjTJF_PZzklffUS8-gc3Mka0zYnEhzhWFUh0pcxpDpjcFA0WWSdUWW_ASVJFXlbYf0Fos_1IVEVs7svfeTq0w8N2H0XXC4gOHIX7qf7heXaLTJFsuMbs3HNoTP0wrwUtQ","e":"AQAB","d":"J7HWnHiu_444ZYugkr0I1uFyMZE4NpwEi7Henk7_ToVmPPsL_7_j-DgDVlVpq--AxrXZwbuTy7gKsyEUblS7MmPdMfxSw09-2XAJ_X_e0ggqLpxZyebZcnvf320KNI6djFA5AvjKC6ZiwfSVmJYp2sGwDrjw1xK5LMfVxBB9O6dfHQIPojEz6dWp5gELreR89PAFNYkO5BXWqprXEjkWqQLInblzCV539fqlWIpfi9RZPRmn5johJu7JlkZe3HxaeQRwJaKdQ1ysYgORLl9pAvY9gys6CcwfohZD9sskeNfA5htn5dRZjDHY828fnfzHgjwjT46zDf-xMqWIySysiQ","p":"5czggQoh46Bqj1lpDYIytzQZ-95kQCEzkySwz-b_qZSgd_d6ZSXxt_SyY1uu3dXKDAeWzTxt_LlRBesIspUwNkHwTbE5nrVYPtgUGehR7CKb2P3UMexOu2qeUkqZH6xATvh-2ZZPmCY7BGi8NPMWp-wBsYo5EnPvwwbq4bmhZp8","q":"3vZioIlNun5ZcoiBAK7kXxN6iF9EC4IFZMM2E8trsuVs7EfWjL6_gPyBctSxn1VixiFl0iLZVDvAs4mmbGj78AUTQ2GQmtVeD0DSGKVm8cfsIqafkIfOIcF-Hl7_eAmdWtNNpx5gD3I43ZyHn7ZfD5w6g57qMQTXLjSGYJ8JKCs","dp":"EVe-8b4kBJvMrvjedsiGr1DdTSbhhf17ePVh6q7SSKgQ3DzvHccZUPrEo779mXxS_UltVhvjaRlLRhkQ1PlxZAbh7dscMCAbgtKn4bSoyhtqi5vMceAVqQtI24kJuVw0lkEmwaEYbLEl7xVAbvaRlSa4kf-OgxgA1kUlYNezmJ0","dq":"vcZtMENt-3yr2cbCNril_R7xPr4HhtwGhzt4_eQ5KS4KRhrnTSjWi41hCUJsZTgiOI4YwoGTBhVN8gMJumCpgCRxvvp-QKu3wbfkm8G9G7KVFPFKA5T0KNsu497sB1n3q2ULRWGfVcZdDJO9BH4P75OEYp-SqmJ6XQOsFPWIr70","qi":"wOZjIsgisv-8anVZ2ZQR-AORbl7BQjxtbtZY9WfQ08a3TJAlWG4xyGN4NIRfQ2jVlc-EpQ30umpGU86ra7wif2w0-DoGP-aYH_1Wg9NOeB2mZ8xqhZogpCKzlhdSkim6HVOIYUl6wqst-aQTjk6wBTTIoJI8_uk5MK3HtK5trNM"}`

	// header
	HeaderNonce = "X-Admin-Nonce"
	HeaderEpoch = "X-Admin-Epoch"

	// lifetime
	adminTokenLifetime = 30
)

func VerifyAdminToken(accessToken string, r *http.Request) error {
	parsedPubJwk := jose.JSONWebKey{}
	if err := parsedPubJwk.UnmarshalJSON([]byte(publicKey)); err != nil {
		return err
	}

	parsedSig, err := jose.ParseSigned(accessToken)
	if err != nil {
		return err
	}

	// verify signature
	expected, err := parsedSig.Verify(parsedPubJwk.Key)
	if err != nil {
		return err
	}

	// verify if payload is correct
	if string(expected) != string(getSignMsg(r)) {
		return fmt.Errorf("signature payload mispatch (expected = %s, received = %s)", expected, getSignMsg(r))
	}

	// check if date is correct
	receivedEpoch, err := strconv.ParseInt(r.Header.Get(HeaderEpoch), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse epoch header: %v", err)
	}

	diff := receivedEpoch - time.Now().Unix()
	if math.Abs(float64(diff)) > adminTokenLifetime {
		return fmt.Errorf("date header not within range (diff = %ds)", diff)
	}

	return nil
}

// AuthorizeAdmin adds 'Authorization' header for internal service to service calls
// Header = Bearer jws(HeaderNonce|HeaderEpoch)
// jws calculates the json web signature of the given message (which is constructed as described above)
// HeaderNonce is a random nonce, HeaderEpoch is the unix epoch in seconds (both are explicitly set by this function)
func AuthorizeAdmin(r *http.Request) *http.Request {

	// explicitly set headers
	r.Header.Set(HeaderNonce, fmt.Sprintf("%d", rand.Int63()))
	r.Header.Set(HeaderEpoch, fmt.Sprintf("%d", time.Now().Unix()))

	// add signature as header
	r.Header.Set(HeaderAuthorization, fmt.Sprintf("%s %s", TokenTypeAdmin, sign(getSignMsg(r))))

	return r
}

// helpers
func getSignMsg(r *http.Request) []byte {
	return []byte(r.Header.Get(HeaderNonce) + r.Header.Get(HeaderEpoch))
}

func sign(data []byte) string {
	parsedJwk := jose.JSONWebKey{}
	if err := parsedJwk.UnmarshalJSON([]byte(privateKey)); err != nil {
		log.Printf("failed to create signature: %v", err)
		return ""
	}

	signingKey := jose.SigningKey{
		Key:       parsedJwk.Key,
		Algorithm: jose.RS256,
	}

	signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		log.Printf("failed to create signature: %v", err)
		return ""
	}

	sig, err := signer.Sign(data)
	if err != nil {
		log.Printf("failed to create signature: %v", err)
		return ""
	}

	sigText, err := sig.CompactSerialize()
	if err != nil {
		log.Printf("failed to create signature: %v", err)
		return ""
	}
	return sigText
}
