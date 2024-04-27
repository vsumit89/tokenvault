# Token Vault

`tokenvault` is a Go package for managing tokens with configurable generation functions and expiration durations. It provides a thread-safe `TokenManager` struct for generating, storing, and retrieving tokens in a concise and efficient manner.

## Features

- **TokenManager**: A struct that encapsulates the logic for generating, storing, and retrieving tokens.
- **Configurable Token Generation**: Supports plugging in custom token generation functions (`TokenGeneratorFunc`) based on application requirements.
- **Configurable Token Duration**: Allows setting the desired duration for token expiration.
- **Thread-Safety**: Utilizes atomic operations to ensure safe access and updates to tokens across multiple goroutines.
- **Blocking Mechanism**: Ensures the first token is generated before allowing access, preventing race conditions.
