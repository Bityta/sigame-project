package com.sigame.lobby.service.batch

import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.grpc.PackServiceClient
import com.sigame.lobby.grpc.UserInfo
import com.sigame.lobby.grpc.PackInfo
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import mu.KotlinLogging
import org.springframework.stereotype.Service
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Service
class BatchOperationService(
    private val authServiceClient: AuthServiceClient,
    private val packServiceClient: PackServiceClient
) {
    
        suspend fun getUserInfoBatch(userIds: List<UUID>): Map<UUID, UserInfo?> = coroutineScope {
        val uniqueIds = userIds.distinct()
        
        logger.debug { "Fetching user info for ${uniqueIds.size} users" }
        
        val results = uniqueIds.map { userId ->
            async {
                userId to authServiceClient.getUserInfo(userId)
            }
        }.map { it.await() }
        
        results.toMap()
    }
    
        suspend fun getPackInfoBatch(packIds: List<UUID>): Map<UUID, PackInfo?> = coroutineScope {
        val uniqueIds = packIds.distinct()
        
        logger.debug { "Fetching pack info for ${uniqueIds.size} packs" }
        
        val results = uniqueIds.map { packId ->
            async {
                packId to packServiceClient.getPackInfo(packId)
            }
        }.map { it.await() }
        
        results.toMap()
    }
}

