package disk_usage

import "testing"

func TestRowToUsage(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}

	usage := rowToUsage("/dev/sda6                   19998104  7083692  12914412  36% /")

	device := "/dev/sda6"
	assertEqual(usage.device, device)

	totalBlocks := 19998104
	assertEqual(usage.total_blocks, totalBlocks)

	availableBlocks := 12914412
	assertEqual(usage.available_blocks, availableBlocks)

	usedBlocks := 7083692
	assertEqual(usage.used_blocks, usedBlocks)

	blockPct := 36
	assertEqual(usage.BlockPct(), blockPct)
}

func TestAddInodeToUsage(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}

	blockLine := "/dev/sda6                   19998104  7083692  12914412  36% /"
	usage := rowToUsage(blockLine)
	usageWithInode := rowToUsage(blockLine)
	inodeLine := "/dev/sda6    124496   16308   108188   14% /"
	addInodeToUsage(inodeLine, usageWithInode)

	assertEqual(usageWithInode.device, usage.device)

	assertEqual(usageWithInode.total_blocks, usage.total_blocks)

	assertEqual(usageWithInode.available_blocks, usage.available_blocks)

	assertEqual(usageWithInode.used_blocks, usage.used_blocks)

	assertEqual(usageWithInode.BlockPct(), usage.BlockPct())

	assertEqual(usage.total_inodes, 0)

	assertEqual(usage.available_inodes, 0)

	assertEqual(usage.used_inodes, 0)

	assertEqual(usage.InodePct(), 0)

	totalInodes := 124496
	assertEqual(usageWithInode.total_inodes, totalInodes)

	availableInodes := 108188
	assertEqual(usageWithInode.available_inodes, availableInodes)

	usedInodes := 16308
	assertEqual(usageWithInode.used_inodes, usedInodes)

	inodesPct := 14
	assertEqual(usageWithInode.InodePct(), inodesPct)
}
