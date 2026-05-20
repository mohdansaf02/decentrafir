// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title FIRManagement
 * @notice Tamper-proof FIR hash storage with role-based access control
 */
contract FIRManagement {
    enum FIRStatus {
        Pending,
        UnderReview,
        Approved,
        Rejected,
        Closed
    }

    struct FIR {
        string firId;
        address citizen;
        bytes32 firHash;
        uint256 timestamp;
        FIRStatus status;
        address assignedOfficer;
        bool exists;
    }

    address public admin;
    mapping(address => bool) public policeOfficers;
    mapping(string => FIR) private firs;
    string[] private firIds;

    event FIRCreated(
        string indexed firId,
        address indexed citizen,
        bytes32 firHash,
        uint256 timestamp
    );
    event FIRUpdated(
        string indexed firId,
        FIRStatus status,
        address assignedOfficer
    );
    event FIRApproved(string indexed firId, address indexed officer);
    event OfficerAdded(address indexed officer);
    event OfficerRemoved(address indexed officer);

    modifier onlyAdmin() {
        require(msg.sender == admin, "Not admin");
        _;
    }

    modifier onlyPolice() {
        require(policeOfficers[msg.sender], "Not police officer");
        _;
    }

    modifier firExists(string memory firId) {
        require(firs[firId].exists, "FIR does not exist");
        _;
    }

    constructor() {
        admin = msg.sender;
        policeOfficers[msg.sender] = true;
    }

    /// @notice Register a new FIR on-chain with metadata hash
    function createFIR(
        string calldata firId,
        bytes32 firHash
    ) external returns (bool) {
        require(!firs[firId].exists, "FIR already exists");
        require(bytes(firId).length > 0, "Invalid FIR ID");

        firs[firId] = FIR({
            firId: firId,
            citizen: msg.sender,
            firHash: firHash,
            timestamp: block.timestamp,
            status: FIRStatus.Pending,
            assignedOfficer: address(0),
            exists: true
        });
        firIds.push(firId);

        emit FIRCreated(firId, msg.sender, firHash, block.timestamp);
        return true;
    }

    /// @notice Update FIR status; police only
    function updateFIRStatus(
        string calldata firId,
        FIRStatus newStatus
    ) external onlyPolice firExists(firId) returns (bool) {
        FIR storage fir = firs[firId];
        fir.status = newStatus;

        if (newStatus == FIRStatus.Approved) {
            emit FIRApproved(firId, msg.sender);
        }

        emit FIRUpdated(firId, newStatus, fir.assignedOfficer);
        return true;
    }

    /// @notice Assign investigating officer to FIR
    function assignOfficer(
        string calldata firId,
        address officer
    ) external onlyPolice firExists(firId) returns (bool) {
        require(policeOfficers[officer], "Invalid officer");
        firs[firId].assignedOfficer = officer;
        emit FIRUpdated(firId, firs[firId].status, officer);
        return true;
    }

    /// @notice Retrieve FIR details
    function getFIR(
        string calldata firId
    )
        external
        view
        firExists(firId)
        returns (
            string memory id,
            address citizen,
            bytes32 firHash,
            uint256 timestamp,
            FIRStatus status,
            address assignedOfficer
        )
    {
        FIR storage fir = firs[firId];
        return (
            fir.firId,
            fir.citizen,
            fir.firHash,
            fir.timestamp,
            fir.status,
            fir.assignedOfficer
        );
    }

    /// @notice Admin adds police officer
    function addOfficer(address officer) external onlyAdmin {
        require(officer != address(0), "Invalid address");
        policeOfficers[officer] = true;
        emit OfficerAdded(officer);
    }

    /// @notice Admin removes police officer
    function removeOfficer(address officer) external onlyAdmin {
        policeOfficers[officer] = false;
        emit OfficerRemoved(officer);
    }

    /// @notice Total FIR count on chain
    function getFIRCount() external view returns (uint256) {
        return firIds.length;
    }
}
