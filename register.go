package workr

// Register a job
func (wrkr *Workr) Register(jj Job) {
	if wrkr.jobs == nil {
		wrkr.jobs = make(map[string]Job)
	}

	wrkr.jobs[jj.Name()] = jj
}
