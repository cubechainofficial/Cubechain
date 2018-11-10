# CubeChain

27 blocks into a cube, the new blockchain platform

Features of CubeChain

## 1. Cubing

Cubing refers to a technology for making a cube by combines 27 blocks into a cube.
The 24 blocks that record the transaction book and three special blocks are combined to create one cube.
The special blocks are basically composed of three blocks but the number of data blocks can be adjusted and the ratio can be adjusted.
As the data block is generated, the cubing is performed and the generated cube creates another hash value.
Since the hash value of the block and the hash value of the cube are generated due to the cube, a data system which can be doubly verified can be constructed.

## 2. Indexing Block

Indexing block is a block indexing whole data concisely and intensifies its search function.
It is a data block that summarizes an E-Wallet of whole transactions by addresses and cube height(block height of existing blockchain).
Thus, an indexing block makes it even faster to find specific data from an EWallet.
Now it is implemented at a high speed when providing a list output history of electronic wallet or API.

## 3. Statistics Block

It is a block that organizes statistical values of whole blocks and expedites the system.
To find a certain wallet, you have to search whole blocks of E-Wallets, make lists of balance, and you have to look for EWallets containing more than 5,000 QUB.
However, if the data of the POS is collected in the statistics block, the search process that is to be repeated every time can be efficiently reduced.
It will be very efficient to search data if you collect frequently used data such as a list of top 1,000 of E-Wallets or a list of E-Wallets with more than 100 times of transfers.
Consequently, the API of the corresponding application service can be implemented at a high speed.

## 4. Escrow Block

An Escrow Block records double approvals data.
A Double authorization data system is a method that allows transactions to be made only after double authorization is issued to authorized traders to issue an approved encryption key during the approval process when using a common blockchain. General data is recorded as one of 24 data. However, escrow data are kept separately. When double approvals are made, they have recorded general data again. Although transactions are made, an escrow block records double approval of cryptocurrencies which are from E-Wallets and unavailable to use immediately.
In addition to an escrow account whose money is held by the third party, it has also a protection function for transactions between the parties.
It is an escrow function based on a blockchain.
Therefore, it can build secure trading system not only for online shopping malls and open markets but also for direct transactions between individuals.
The double approvals method can be implemented in many ways such as sender approval method, recipient approval method, both sides approval, automatic approval method after a specific period, and so on. Escrow blocks can be used by the owner to protect data through passwords, which means only users who know the password through encryption can view the data, rather than using the data in an open format.

## 5. POH (Proof of POW+POS hybrid)

Cube Chain adopts POW and POS hybrid method.
In the beginning, the ratio of POW to POS is 7:3. As time goes, the ratio of POS increases and finally it will be maintained only by POS.
It increases POW in the beginning stage to make the network stable and raises POS gradually to reduce network resources and power waste.
In spite of time-consuming disadvantages for payment when POW and POS are used together, the statistics block of Cube Chain will drastically reduce inefficiency to repeatedly calculate each time.
