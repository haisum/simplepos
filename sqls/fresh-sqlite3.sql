CREATE TABLE Item(
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Name text,
  Attributes text,
  SellPrice INT,
  BuyPrice INT,
  Stock INT,
  Comments text
);
CREATE INDEX Item_Stock ON Item(Stock);

CREATE TABLE Sell(
  ID INTEGER PRIMARY KEY AUTOINCREMENT ,
  Seller TEXT
);

CREATE TABLE SellItem(
  SellID INTEGER REFERENCES Sell(ID),
  ItemID INTEGER REFERENCES Item(ID),
  SellPrice INTEGER,
  Discount INTEGER,
  Quantity INTEGER
);
CREATE INDEX SellItem_SellID ON Sell(ID);
CREATE INDEX SellItem_ItemID ON Item(ID);

CREATE TABLE Buy(
  ID INTEGER PRIMARY KEY AUTOINCREMENT ,
  Buyer TEXT
);

CREATE TABLE BuyItem(
  BuyID INTEGER REFERENCES Buy(ID),
  ItemID INTEGER REFERENCES Item(ID),
  BuyPrice INTEGER,
  Quantity INTEGER
);
CREATE INDEX BuyItem_BuyID ON Buy(ID);
CREATE INDEX BuyItem_ItemID ON Item(ID);

CREATE TABLE Settings(
  Name TEXT,
  Value TEXT
);
CREATE INDEX Settings_Name on Settings(Name);