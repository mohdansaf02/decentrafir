const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("FIRManagement", function () {
  let contract, admin, citizen, officer;

  beforeEach(async function () {
    [admin, citizen, officer] = await ethers.getSigners();
    const FIRManagement = await ethers.getContractFactory("FIRManagement");
    contract = await FIRManagement.deploy();
    await contract.waitForDeployment();
    await contract.connect(admin).addOfficer(officer.address);
  });

  it("creates FIR", async function () {
    const hash = ethers.keccak256(ethers.toUtf8Bytes("fir-metadata"));
    await expect(contract.connect(citizen).createFIR("FIR-001", hash))
      .to.emit(contract, "FIRCreated");

    const fir = await contract.getFIR("FIR-001");
    expect(fir[0]).to.equal("FIR-001");
    expect(fir[1]).to.equal(citizen.address);
  });

  it("updates status by police", async function () {
    const hash = ethers.keccak256(ethers.toUtf8Bytes("data"));
    await contract.connect(citizen).createFIR("FIR-002", hash);
    await contract.connect(officer).updateFIRStatus("FIR-002", 2); // Approved
    const fir = await contract.getFIR("FIR-002");
    expect(fir[4]).to.equal(2n);
  });

  it("assigns officer", async function () {
    const hash = ethers.keccak256(ethers.toUtf8Bytes("data"));
    await contract.connect(citizen).createFIR("FIR-003", hash);
    await contract.connect(officer).assignOfficer("FIR-003", officer.address);
    const fir = await contract.getFIR("FIR-003");
    expect(fir[5]).to.equal(officer.address);
  });

  it("rejects duplicate FIR", async function () {
    const hash = ethers.keccak256(ethers.toUtf8Bytes("data"));
    await contract.connect(citizen).createFIR("FIR-004", hash);
    await expect(
      contract.connect(citizen).createFIR("FIR-004", hash)
    ).to.be.revertedWith("FIR already exists");
  });
});
