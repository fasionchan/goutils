/*
 * Author: fasion
 * Created time: 2023-03-22 17:20:15
 * Last Modified by: fasion
 * Last Modified time: 2023-03-22 17:21:40
 */

package jobutils

func WaitByChan(f func()) chan struct{} {
	c := make(chan struct{})
	go func() {
		f()
		c <- struct{}{}
	}()
	return c
}
