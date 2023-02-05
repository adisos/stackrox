// Package netpol provides primitives for command 'roxctl analyze netpol'
package netpol

func (cmd *analyzeNetpolCommand) printNetpolAnalysis(analysisRes string) {
	cmd.env.Logger().PrintfLn(analysisRes)
}
