class AppConfig {
  static const apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );
  static const contractAddress = String.fromEnvironment(
    'CONTRACT_ADDRESS',
    defaultValue: '',
  );
  static const rpcUrl = String.fromEnvironment(
    'RPC_URL',
    defaultValue: 'https://rpc-mumbai.maticvigil.com',
  );
  static const mumbaiChainId = 80001;
}
