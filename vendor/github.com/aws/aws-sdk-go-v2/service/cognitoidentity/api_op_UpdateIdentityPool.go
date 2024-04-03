// Code generated by smithy-go-codegen DO NOT EDIT.

package cognitoidentity

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Updates an identity pool. You must use AWS Developer credentials to call this
// API.
func (c *Client) UpdateIdentityPool(ctx context.Context, params *UpdateIdentityPoolInput, optFns ...func(*Options)) (*UpdateIdentityPoolOutput, error) {
	if params == nil {
		params = &UpdateIdentityPoolInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateIdentityPool", params, optFns, c.addOperationUpdateIdentityPoolMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateIdentityPoolOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// An object representing an Amazon Cognito identity pool.
type UpdateIdentityPoolInput struct {

	// TRUE if the identity pool supports unauthenticated logins.
	//
	// This member is required.
	AllowUnauthenticatedIdentities bool

	// An identity pool ID in the format REGION:GUID.
	//
	// This member is required.
	IdentityPoolId *string

	// A string that you provide.
	//
	// This member is required.
	IdentityPoolName *string

	// Enables or disables the Basic (Classic) authentication flow. For more
	// information, see Identity Pools (Federated Identities) Authentication Flow (https://docs.aws.amazon.com/cognito/latest/developerguide/authentication-flow.html)
	// in the Amazon Cognito Developer Guide.
	AllowClassicFlow *bool

	// A list representing an Amazon Cognito user pool and its client ID.
	CognitoIdentityProviders []types.CognitoIdentityProvider

	// The "domain" by which Cognito will refer to your users.
	DeveloperProviderName *string

	// The tags that are assigned to the identity pool. A tag is a label that you can
	// apply to identity pools to categorize and manage them in different ways, such as
	// by purpose, owner, environment, or other criteria.
	IdentityPoolTags map[string]string

	// The ARNs of the OpenID Connect providers.
	OpenIdConnectProviderARNs []string

	// An array of Amazon Resource Names (ARNs) of the SAML provider for your identity
	// pool.
	SamlProviderARNs []string

	// Optional key:value pairs mapping provider names to provider app IDs.
	SupportedLoginProviders map[string]string

	noSmithyDocumentSerde
}

// An object representing an Amazon Cognito identity pool.
type UpdateIdentityPoolOutput struct {

	// TRUE if the identity pool supports unauthenticated logins.
	//
	// This member is required.
	AllowUnauthenticatedIdentities bool

	// An identity pool ID in the format REGION:GUID.
	//
	// This member is required.
	IdentityPoolId *string

	// A string that you provide.
	//
	// This member is required.
	IdentityPoolName *string

	// Enables or disables the Basic (Classic) authentication flow. For more
	// information, see Identity Pools (Federated Identities) Authentication Flow (https://docs.aws.amazon.com/cognito/latest/developerguide/authentication-flow.html)
	// in the Amazon Cognito Developer Guide.
	AllowClassicFlow *bool

	// A list representing an Amazon Cognito user pool and its client ID.
	CognitoIdentityProviders []types.CognitoIdentityProvider

	// The "domain" by which Cognito will refer to your users.
	DeveloperProviderName *string

	// The tags that are assigned to the identity pool. A tag is a label that you can
	// apply to identity pools to categorize and manage them in different ways, such as
	// by purpose, owner, environment, or other criteria.
	IdentityPoolTags map[string]string

	// The ARNs of the OpenID Connect providers.
	OpenIdConnectProviderARNs []string

	// An array of Amazon Resource Names (ARNs) of the SAML provider for your identity
	// pool.
	SamlProviderARNs []string

	// Optional key:value pairs mapping provider names to provider app IDs.
	SupportedLoginProviders map[string]string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateIdentityPoolMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpUpdateIdentityPool{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpUpdateIdentityPool{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateIdentityPool"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpUpdateIdentityPoolValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateIdentityPool(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opUpdateIdentityPool(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateIdentityPool",
	}
}
