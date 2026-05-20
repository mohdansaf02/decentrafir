const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying with account:", deployer.address);

  const FIRManagement = await hre.ethers.getContractFactory("FIRManagement");
  const contract = await FIRManagement.deploy();
  await contract.waitForDeployment();

  const address = await contract.getAddress();
  console.log("FIRManagement deployed to:", address);

  const outDir = path.join(__dirname, "..", "..", "backend", "config");
  fs.mkdirSync(outDir, { recursive: true });
  fs.writeFileSync(
    path.join(outDir, "contract.json"),
    JSON.stringify(
      {
        address,
        network: hre.network.name,
        chainId: (await hre.ethers.provider.getNetwork()).chainId.toString(),
        deployedAt: new Date().toISOString(),
      },
      null,
      2
    )
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
