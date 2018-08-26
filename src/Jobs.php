<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\Jobs;

use Spiral\Core\Container\SingletonInterface;
use Spiral\Goridge\RPC;
use Spiral\Jobs\Configs\JobsConfig;
use Spiral\Jobs\Exceptions\JobException;
use Spiral\RoadRunner\Exceptions\RoadRunnerException;

class Jobs implements JobsInterface, SingletonInterface
{
    const RR_SERVICE = 'jobs';

    /** @var JobsConfig */
    private $config;

    /** @var RPC */
    private $rpc;

    /**
     * @param JobsConfig $config
     * @param RPC        $rpc
     */
    public function __construct(JobsConfig $config, RPC $rpc)
    {
        $this->config = $config;
        $this->rpc = $rpc;
    }

    /**
     * @inheritdoc
     */
    public function push(JobInterface $job, Options $options = null): string
    {
        try {
            if (empty($options)) {
                $options = new Options();
            }

            return $this->rpc->call(self::RR_SERVICE . '.Push', [
                'job'      => get_class($job),
                'pipeline' => $options->getPipeline() ?? $this->config->jobPipeline(get_class($job)),
                'payload'  => $job,
                'options'  => $options
            ]);
        } catch (RoadRunnerException|\Throwable $e) {
            throw new JobException($e->getMessage(), $e->getCode(), $e);
        }
    }
}