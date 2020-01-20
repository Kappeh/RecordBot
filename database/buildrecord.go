package database

// Build gets the build of the build record
func (b BuildRecord) Build() (Build, error) {
	return Build{}, nil
}

// Record gets the record of the build record
func (b BuildRecord) Record() (Record, error) {
	return Record{}, nil
}

// FirstJointBuildRecord gets the first (chronologically) joint build record
func (b BuildRecord) FirstJointBuildRecord() (BuildRecord, error) {
	return BuildRecord{}, nil
}

// JointBuildRecords gets all joint build records
func (b BuildRecord) JointBuildRecords() ([]BuildRecord, error) {
	return nil, nil
}
