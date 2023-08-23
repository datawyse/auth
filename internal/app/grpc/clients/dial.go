package clients

// func dial(config *internal.AppConfig) (*grpc.ClientConn, error) {
// 	// ta, err := domain.NewKeycloakTokenAuth(config.KeycloakServer, config)
// 	// if err != nil {
// 	// 	return nil, errors.Wrapf(err, "NewKeycloakTokenAuth")
// 	// }
//
// 	opts := []grpc.DialOption{
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
//
// 		// In addition to the following grpc.DialOption, callers may also use
// 		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
// 		// itself.
// 		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
// 		// grpc.WithPerRPCCredentials(ta),
//
// 		// oauth.TokenSource requires the configuration of transport credentials.
// 		// grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
// 		// 	InsecureSkipVerify: true,
// 		// })),
// 	}
//
// 	conn, err := grpc.Dial(config.LoggingGrpcService, opts...)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "grpc.Dial")
// 	}
//
// 	return conn, nil
// }
