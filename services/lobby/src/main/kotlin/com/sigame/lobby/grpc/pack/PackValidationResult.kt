package com.sigame.lobby.grpc.pack

data class PackValidationResult(
    val exists: Boolean,
    val isOwner: Boolean,
    val status: String,
    val error: String?
) {
    val isApproved: Boolean get() = status == "approved"
}

