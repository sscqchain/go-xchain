total_provision = 3250000000
initial_provision = 325000000
mine_provision = total_provision - initial_provision

initial_reward =278.253424658
years = 50
blocktime = 6
blks_per_day = 24 * 3600/blocktime

time_per_year = 365 * 24 * 3600
blocks_per_year = time_per_year/blocktime

if __name__ == "__main__":
    accumulated_provision = 0
    reward = initial_reward
    for i in range(50):
        accumulated_provision += reward * blocks_per_year
        reward = reward/2.0
    print("expected mining provision: ",mine_provision)
    print("accumulated_provision: ", accumulated_provision)
    assert mine_provision == int(accumulated_provision)
    print("initial day provision:", blks_per_day*initial_reward)