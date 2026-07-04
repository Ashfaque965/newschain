// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title NewsRegistry
/// @notice Anchors news article content hashes on-chain to provide
///         tamper-evident, timestamped proof of authorship and publication.
///         The full article body lives off-chain (IPFS / Go backend + SQLite);
///         only a keccak256 hash + metadata pointer is stored here.
contract NewsRegistry {
    struct Article {
        address author;         // wallet that published the article
        bytes32 contentHash;    // keccak256 hash of canonical article body
        string  ipfsCID;        // IPFS CID pointing to full article content
        uint256 publishedAt;    // block timestamp of first publication
        uint256 editCount;      // number of times content has been amended
        bool    retracted;      // true if author/moderator retracted it
    }

    /// articleId => Article
    mapping(bytes32 => Article) private articles;

    /// author => list of articleIds they published (for quick lookup/UI)
    mapping(address => bytes32[]) private authorArticles;

    /// Moderators who can flag disputed articles (does not alter content)
    mapping(address => bool) public moderators;

    address public owner;

    event ArticlePublished(
        bytes32 indexed articleId,
        address indexed author,
        bytes32 contentHash,
        string ipfsCID,
        uint256 timestamp
    );

    event ArticleAmended(
        bytes32 indexed articleId,
        bytes32 newContentHash,
        string newIpfsCID,
        uint256 timestamp
    );

    event ArticleRetracted(bytes32 indexed articleId, uint256 timestamp);
    event ArticleDisputed(bytes32 indexed articleId, address indexed moderator, string reason);

    modifier onlyOwner() {
        require(msg.sender == owner, "NewsRegistry: not owner");
        _;
    }

    modifier onlyAuthor(bytes32 articleId) {
        require(articles[articleId].author == msg.sender, "NewsRegistry: not author");
        _;
    }

    modifier onlyModerator() {
        require(moderators[msg.sender] || msg.sender == owner, "NewsRegistry: not moderator");
        _;
    }

    constructor() {
        owner = msg.sender;
        moderators[msg.sender] = true;
    }

    /// @notice Publish a new article's fingerprint on-chain.
    /// @param articleId Deterministic ID, e.g. keccak256(abi.encodePacked(author, slug))
    /// @param contentHash keccak256 hash of the canonical article body (computed off-chain)
    /// @param ipfsCID IPFS content identifier where the full article is stored
    function publishArticle(
        bytes32 articleId,
        bytes32 contentHash,
        string calldata ipfsCID
    ) external {
        require(articles[articleId].publishedAt == 0, "NewsRegistry: article exists");
        require(contentHash != bytes32(0), "NewsRegistry: empty hash");

        articles[articleId] = Article({
            author: msg.sender,
            contentHash: contentHash,
            ipfsCID: ipfsCID,
            publishedAt: block.timestamp,
            editCount: 0,
            retracted: false
        });

        authorArticles[msg.sender].push(articleId);

        emit ArticlePublished(articleId, msg.sender, contentHash, ipfsCID, block.timestamp);
    }

    /// @notice Amend an existing article (e.g. correction). Only the original author may amend.
    function amendArticle(
        bytes32 articleId,
        bytes32 newContentHash,
        string calldata newIpfsCID
    ) external onlyAuthor(articleId) {
        require(articles[articleId].publishedAt != 0, "NewsRegistry: not found");
        require(!articles[articleId].retracted, "NewsRegistry: retracted");

        Article storage a = articles[articleId];
        a.contentHash = newContentHash;
        a.ipfsCID = newIpfsCID;
        a.editCount += 1;

        emit ArticleAmended(articleId, newContentHash, newIpfsCID, block.timestamp);
    }

    /// @notice Retract an article. Content stays on IPFS but is flagged retracted on-chain.
    function retractArticle(bytes32 articleId) external onlyAuthor(articleId) {
        require(articles[articleId].publishedAt != 0, "NewsRegistry: not found");
        articles[articleId].retracted = true;
        emit ArticleRetracted(articleId, block.timestamp);
    }

    /// @notice Flag an article as disputed (e.g. factual challenge). Does not alter content.
    function disputeArticle(bytes32 articleId, string calldata reason) external onlyModerator {
        require(articles[articleId].publishedAt != 0, "NewsRegistry: not found");
        emit ArticleDisputed(articleId, msg.sender, reason);
    }

    /// @notice Verify that a given content hash matches what's anchored on-chain.
    function verify(bytes32 articleId, bytes32 candidateHash) external view returns (bool) {
        return articles[articleId].contentHash == candidateHash && articles[articleId].publishedAt != 0;
    }

    function getArticle(bytes32 articleId) external view returns (Article memory) {
        return articles[articleId];
    }

    function getAuthorArticles(address author) external view returns (bytes32[] memory) {
        return authorArticles[author];
    }

    function setModerator(address account, bool status) external onlyOwner {
        moderators[account] = status;
    }
}
