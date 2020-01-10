/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"github.com/hyperledger/fabric/common/ledger/util"
)

func TestConstructCheckpointInfoFromBlockFiles(t *testing.T) {
	testPath := "/tmp/tests/fabric/common/ledger/blkstorage/fsblkstorage"
	ledgerid := "testLedger"
	conf := NewConf(testPath, 0)
	blkStoreDir := conf.getLedgerBlockDir(ledgerid)
	env := newTestEnv(t, conf)
	util.CreateDirIfMissing(blkStoreDir)
	defer env.Cleanup()

	// checkpoint constructed on an empty block folder should return CPInfo with isChainEmpty: true
	cpInfo, err := constructCheckpointInfoFromBlockFiles(blkStoreDir)
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, cpInfo, &checkpointInfo{isChainEmpty: true, lastBlockNumber: 0, latestFileChunksize: 0, latestFileChunkSuffixNum: 0})

	w := newTestBlockfileWrapper(env, ledgerid)
	defer w.close()
	blockfileMgr := w.blockfileMgr
	bg, gb := testutil.NewBlockGenerator(t, ledgerid, false)

	// Add a few blocks and verify that cpinfo derived from filesystem should be same as from the blockfile manager
	blockfileMgr.addBlock(gb)
	for _, blk := range bg.NextTestBlocks(3) {
		blockfileMgr.addBlock(blk)
	}
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	// Move the chain to new file and check cpinfo derived from file system
	blockfileMgr.moveToNextFile()
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	// Add a few blocks that would go to new file and verify that cpinfo derived from filesystem should be same as from the blockfile manager
	for _, blk := range bg.NextTestBlocks(3) {
		blockfileMgr.addBlock(blk)
	}
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	// Write a partial block (to simulate a crash) and verify that cpinfo derived from filesystem should be same as from the blockfile manager
	lastTestBlk := bg.NextTestBlocks(1)[0]
	blockBytes, _, err := serializeBlock(lastTestBlk)
	testutil.AssertNoError(t, err, "")
	partialByte := append(proto.EncodeVarint(uint64(len(blockBytes))), blockBytes[len(blockBytes)/2:]...)
	blockfileMgr.currentFileWriter.append(partialByte, true)
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)

	// Close the block storage, drop the index and restart and verify
	cpInfoBeforeClose := blockfileMgr.cpInfo
	w.close()
	env.provider.Close()
	indexFolder := conf.getIndexDir()
	testutil.AssertNoError(t, os.RemoveAll(indexFolder), "")

	env = newTestEnv(t, conf)
	w = newTestBlockfileWrapper(env, ledgerid)
	blockfileMgr = w.blockfileMgr
	testutil.AssertEquals(t, blockfileMgr.cpInfo, cpInfoBeforeClose)

	lastBlkIndexed, err := blockfileMgr.index.getLastBlockIndexed()
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, lastBlkIndexed, uint64(6))

	// Add the last block again after start and check cpinfo again
	testutil.AssertNoError(t, blockfileMgr.addBlock(lastTestBlk), "")
	checkCPInfoFromFile(t, blkStoreDir, blockfileMgr.cpInfo)
}

func checkCPInfoFromFile(t *testing.T, blkStoreDir string, expectedCPInfo *checkpointInfo) {
	cpInfo, err := constructCheckpointInfoFromBlockFiles(blkStoreDir)
	testutil.AssertNoError(t, err, "")
	testutil.AssertEquals(t, cpInfo, expectedCPInfo)
}
