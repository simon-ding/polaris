package engine

import "polaris/log"

func (c *Engine) housekeeping() error {
	log.Infof("start housekeeping tasks...")
	
	if err := c.checkDbScraps(); err != nil {
		return err
	}
	if err := c.checkImageFilesInterity(); err != nil {
		return err
	}
	
	return nil
}

func (c *Engine) checkDbScraps() error {
	//TODO: remove episodes that are not associated with any series
	return nil
}

func (c *Engine) checkImageFilesInterity() error {
	//TODO: download missing image files, remove unused image files
	return nil
}