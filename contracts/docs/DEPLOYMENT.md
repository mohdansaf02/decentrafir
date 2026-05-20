# Smart Contract Deployment Guide

## Local (Hardhat)

```bash
cd contracts
npm install
npx hardhat node          # terminal 1
npm run deploy:local      # terminal 2
```

Copy the deployed address to `backend/.env` as `CONTRACT_ADDRESS`.

## Polygon Mumbai

1. Fund deployer wallet with test MATIC.
2. Copy `.env.example` to `.env` and set `DEPLOYER_PRIVATE_KEY`.
3. Deploy:

```bash
npm run deploy:mumbai
```

4. Verify on [Polygonscan Mumbai](https://mumbai.polygonscan.com/).

## Post-deploy

- Set `CONTRACT_ADDRESS` in backend and Flutter `--dart-define`.
- Add deployer account as contract admin (constructor assigns deployer).
- Register police wallet addresses via `addOfficer()` on-chain.
