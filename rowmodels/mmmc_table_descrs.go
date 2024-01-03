package rowmodels

import DRU "github.com/fbaube/datarepo/utils"

// MmmcTableDescriptors configures the three key tables.
var MmmcTableDescriptors = []DRU.TableDescriptor{
	TableDescriptor_InbatchRow,
	TableDescriptor_ContentityRow,
	TableDescriptor_TopicrefRow,
}
